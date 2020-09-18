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

func TestCRUDOperationsForBackupPolicy(t *testing.T) {
	org, err := tests.GetDefaultOrganization()
	require.NoError(t, err)
	proj, err := tests.GetDefaultProject(org)
	require.NoError(t, err)

	cmd.RootCmd.PersistentPreRun(nil, nil)
	ctx := cmd.ContextWithToken()
	conn := cmd.MustDialAPI()
	defer conn.Close()
	datac := data.NewDataServiceClient(conn)
	region, err := getRegion("aws", org)
	require.NoError(t, err)
	version, err := datac.GetDefaultVersion(ctx, &common.Empty{})
	require.NoError(t, err)
	deplName := "TestCreateDeployment"
	cert, err := getDefaultCertificate(proj)
	require.NoError(t, err)

	args := []string{"create", "deployment", "--name=" + deplName, "--organization-id=" + org, "--project-id=" +
		proj, "--region-id=" + region, "--version=" + version.GetVersion(), "--cacertificate-id=" + cert.GetId(), "--accept"}
	out, err := tests.RunCommand(args)
	require.NoError(t, err)
	deplId, err := tests.GetResourceID(string(out))
	require.NoError(t, err)

	defer func() {
		if _, err := datac.DeleteDeployment(ctx, &common.IDOptions{Id: deplId}); err != nil {
			t.Log(err)
		}
	}()

	testBackupPolicy := "TestBackupPolicy"
	args = []string{"create", "backup", "policy", "--name=" + testBackupPolicy, "--deployment-id=" + deplId, "--schedule-type=hourly", "--every-interval-hours=1"}
	compare := `Success!
Id                            .*
Deleted                       -
Deployment-Id                 ` + deplId + `
Description                   
Name                          ` + testBackupPolicy + `
Upload                        -
Url                           /Organization/` + org + `/Project/` + proj + `/Deployment/` + deplId + `/BackupPolicy/.*
Locked                        -
Paused                        -
Schedule-Type                 Hourly
Retention-Period              
Created-At                    .*
Deleted-At                    
Schedule-Every-Interval-Hours 1
`
	out, err = tests.RunCommand(args)
	require.NoError(t, err)
	assert.True(t, tests.CompareOutput(out, []byte(compare)))

	backupPolicyId, err := tests.GetResourceID(string(out))
	require.NoError(t, err)

	t.Run("get backup policy", func(tt *testing.T) {
		args := []string{"get", "backup", "policy", "--id=" + backupPolicyId}
		compare := `Id                            .*
Deleted                       -
Deployment-Id                 ` + deplId + `
Description                   
Name                          ` + testBackupPolicy + `
Upload                        -
Url                           /Organization/` + org + `/Project/` + proj + `/Deployment/` + deplId + `/BackupPolicy/.*
Locked                        -
Paused                        -
Schedule-Type                 Hourly
Retention-Period              
Created-At                    .*
Deleted-At                    
`
		out, err := tests.RunCommand(args)
		require.NoError(tt, err)
		assert.True(tt, tests.CompareOutput(out, []byte(compare)))
	})

	t.Run("list backup policies", func(tt *testing.T) {
		args := []string{"list", "backup", "policies", "--deployment-id=" + deplId}
		compare := `Id\s+| Deleted\s+| Deployment-Id\s+| Description\s+| Name\s+| Upload\s+| Url\s+| Locked\s+| Paused\s+| Schedule-Type\s+` +
			`| Retention-Period\s+| Created-At\s+| Deleted-At\s+| Schedule-Every-Interval-Hours(\s.*)+` +
			backupPolicyId + ` | -\s+| ` + deplId + `\s+|.*| ` + testBackupPolicy + `\s+| -\s+| /Organization/_support/Project/` + proj +
			`/Deployment/` + deplId + `/BackupPolicy/` + backupPolicyId + ` | -\s+| -\s+| Hourly\s+|.*|.*|.*| 1`
		out, err := tests.RunCommand(args)
		require.NoError(tt, err)
		assert.True(tt, tests.CompareOutput(out, []byte(compare)))
	})

	t.Run("update backup policy", func(tt *testing.T) {
		args := []string{"update", "backup", "policy", "--backup-policy-id=" + backupPolicyId, "--name=NewName", "--retention-period=6"}
		compare := `Updated backup policy!
Id                            ` + backupPolicyId + `
Deleted                       -
Deployment-Id                 ` + deplId + `
Description                   
Name                          NewName
Upload                        -
Url                           /Organization/` + org + `/Project/` + proj + `/Deployment/` + deplId + `/BackupPolicy/` + backupPolicyId + `
Locked                        -
Paused                        -
Schedule-Type                 Hourly
Retention-Period              6h0m0s
Created-At                    .*
Deleted-At                    
`
		out, err := tests.RunCommand(args)
		require.NoError(tt, err)
		assert.True(tt, tests.CompareOutput(out, []byte(compare)))
	})

	t.Run("delete backup policy", func(tt *testing.T) {
		args := []string{"delete", "backup", "policy", "--id=" + backupPolicyId}
		compare := `Deleted backup policy!
`
		out, err := tests.RunCommand(args)
		require.NoError(tt, err)
		assert.True(tt, tests.CompareOutput(out, []byte(compare)))
	})
}
