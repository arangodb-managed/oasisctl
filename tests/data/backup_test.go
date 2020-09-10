//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
//
// Author Gergely Brautigam
//
// +build e2e

package data

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	common "github.com/arangodb-managed/apis/common/v1"
	data "github.com/arangodb-managed/apis/data/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/tests"
)

func TestCreateBackup(t *testing.T) {
	org, err := tests.GetDefaultOrganization()
	require.NoError(t, err)
	proj, err := tests.GetDefaultProject(org)
	require.NoError(t, err)

	cmd.RootCmd.PersistentPreRun(nil, nil)
	ctx := cmd.ContextWithToken()
	conn := cmd.MustDialAPI()
	datac := data.NewDataServiceClient(conn)
	region, err := getRegion("aws", org)
	require.NoError(t, err)
	version, err := datac.GetDefaultVersion(ctx, &common.Empty{})
	require.NoError(t, err)
	deplName := "TestCreateDeployment"
	cert, err := getDefaultCertificate(proj)
	require.NoError(t, err)
	tandc, err := getTermsAndConditions(org)
	require.NoError(t, err)
	depl, err := datac.CreateDeployment(ctx, &data.Deployment{
		Name:                         deplName,
		RegionId:                     region,
		Version:                      version.GetVersion(),
		ProjectId:                    proj,
		Certificates:                 &data.Deployment_CertificateSpec{CaCertificateId: cert.GetId()},
		AcceptedTermsAndConditionsId: tandc,
		Model: &data.Deployment_ModelSpec{
			Model:        data.ModelOneShard,
			NodeSizeId:   "c4-a4",
			NodeCount:    3,
			NodeDiskSize: 10,
		},
	})
	require.NoError(t, err)

	defer func() {
		if _, err := datac.DeleteDeployment(ctx, &common.IDOptions{Id: depl.GetId()}); err != nil {
			t.Log(err)
		}
	}()

	testBackup := "TestBackup"
	args := []string{"create", "backup", "--name=" + testBackup, "--deployment-id=" + depl.GetId()}
	compare := `Success!
Id               .*
Backup-Policy-Id 
Deleted          false
Deployment-Id    ` + depl.GetId() + `
Description      
Name             ` + testBackup + `
Upload           false
Url              .*
Auto-Deleted-At  5 hours from now
Created-At       now
Deleted-At       
Dbservers        3
`
	out, err := tests.RunCommand(args)
	require.NoError(t, err)
	assert.True(t, tests.CompareOutput(out, []byte(compare)))

	backupId, err := tests.GetResourceID(string(out))
	require.NoError(t, err)

	t.Run("get backup", func(tt *testing.T) {
		args := []string{"get", "backup", "--id=" + backupId}
		compare := `Id               ` + backupId + `
Backup-Policy-Id 
Deleted          false
Deployment-Id    ` + depl.GetId() + `
Description      
Name             ` + testBackup + `
Upload           false
Url              .*
Auto-Deleted-At  5 hours from now
Created-At       now
Deleted-At       
Dbservers        3
`
		out, err := tests.RunCommand(args)
		require.NoError(tt, err)
		assert.True(tt, tests.CompareOutput(out, []byte(compare)))
	})

	t.Run("list backups", func(tt *testing.T) {
		args := []string{"list", "backups", "--deployment-id=" + depl.GetId()}
		compare := `Id\s+| Backup-Policy-Id | Deleted | Deployment-Id\s+| Description | Name\s+| Upload | Url\s+| State | Db-Servers | Uploaded | Auto-Deleted-At\s+| Created-At | Deleted-At(\s.*)*` + backupId + ` |.*| false\s+| ` + depl.GetId() + ` |.*| TestBackup | false  | /Organization/` + org + `/Project/` + proj + `/Deployment/` + depl.GetId() + `/Backup/` + backupId + ` |.*| 3\s+| false\s+| .* |.*| `
		out, err := tests.RunCommand(args)
		require.NoError(tt, err)
		assert.True(tt, tests.CompareOutput(out, []byte(compare)))
	})

	t.Run("update backup", func(tt *testing.T) {
		args := []string{"update", "backup", "--backup-id=" + backupId, "--name=NewName", "--auto-deleted-at=5"}
		compare := `Updated backup!
Id               ` + backupId + `
Backup-Policy-Id 
Deleted          false
Deployment-Id    ` + depl.GetId() + `
Description      
Name             NewName
Upload           false
Url              .*
Auto-Deleted-At  4 hours from now
Created-At       .*
Deleted-At       
Dbservers        3
`
		out, err := tests.RunCommand(args)
		require.NoError(tt, err)
		assert.True(tt, tests.CompareOutput(out, []byte(compare)))
	})

	t.Run("delete backup", func(tt *testing.T) {
		args := []string{"delete", "backup", "--id=" + backupId}
		compare := `Deleted backup!
`
		out, err := tests.RunCommand(args)
		require.NoError(tt, err)
		assert.True(tt, tests.CompareOutput(out, []byte(compare)))
	})
}
