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

func TestCRUDOperationsForDeployment(t *testing.T) {
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

	compare := `Success!
Id                      .*
Name                    ` + deplName + `
Description             
Region                  .*
Version                 .*
Ipallowlist             
Url                     .*
Paused                  -
Created-At              now
Deleted-At              -
Expires-At              -
Ready                   -
Bootstrapped            -
Created                 -
Upgrading               -
Coordinators            3
Coordinator-Memory-Size 1GB
Dbservers               3
Dbserver-Memory-Size    3GB
Dbserver-Disk-Size      10GB
Bootstrapped-At         -
Endpoint-Url            
Root-Password           \*\*\* use \'--show-root-password\' to expose \*\*\*
Model                   oneshard
Is-Clone                false
Clone-Backup-Id         
Node-Count              3
Node-Disk-Size          10GB
Node-Size-Id            c4-a4
`
	out, err := tests.RunCommand(args)
	require.NoError(t, err)
	assert.True(t, tests.CompareOutput(out, []byte(compare)))
	deplId, err := tests.GetResourceID(string(out))
	require.NoError(t, err)

	defer func() {
		if _, err := datac.DeleteDeployment(ctx, &common.IDOptions{Id: deplId}); err != nil {
			t.Log(err)
		}
	}()

	t.Run("get a deployment", func(tt *testing.T) {
		args := []string{"get", "deployment", "--deployment-id=" + deplId}
		compare := `Id                      ` + deplId + `
Name                    ` + deplName + `
Description             
Region                  .*
Version                 .*
Ipallowlist             
Url                     .*
Paused                  -
Created-At              .*
Deleted-At              -
Expires-At              -
Ready                   -
Bootstrapped            -
Created                 -
Upgrading               -
Coordinators            3
Coordinator-Memory-Size 1GB
Dbservers               3
Dbserver-Memory-Size    3GB
Dbserver-Disk-Size      10GB
Bootstrapped-At         -
Endpoint-Url            .*
Root-Password           \*\*\* use \'--show-root-password\' to expose \*\*\*
Model                   oneshard
Is-Clone                false
Clone-Backup-Id         
Node-Count              3
Node-Disk-Size          10GB
Node-Size-Id            c4-a4
`
		out, err := tests.RunCommand(args)
		assert.NoError(tt, err)
		assert.True(tt, tests.CompareOutput(out, []byte(compare)))
	})

	t.Run("list deployments", func(tt *testing.T) {
		args := []string{"list", "deployments", "--project-id=" + proj}
		compare := `Id\s+| Name\s+| Description\s+| Region\s+| Version | Ipallowlist | Url\s+| Paused | Created-At\s+| Model\s+| Node-Count | Node-Disk-Size | Node-Size-Id(\s.*)*` +
			deplId + ` | ` + deplName + ` |\s+| ` + region + ` | ` + version.GetVersion() + `\s+|\s+| .* | -\s+| .* | oneshard | 3\s+| 10GB\s+| c4-a4
`
		out, err := tests.RunCommand(args)
		assert.NoError(tt, err)
		assert.True(tt, tests.CompareOutput(out, []byte(compare)))
	})

	t.Run("update a deployment", func(tt *testing.T) {
		args := []string{"update", "deployment", "--deployment-id=" + deplId, "--name=NewName"}
		compare := `Updated deployment!
Id                      ` + deplId + `
Name                    NewName
Description             
Region                  .*
Version                 .*
Ipallowlist             
Url                     .*
Paused                  -
Created-At              .*
Deleted-At              -
Expires-At              -
Ready                   -
Bootstrapped            -
Created                 -
Upgrading               -
Coordinators            3
Coordinator-Memory-Size 1GB
Dbservers               3
Dbserver-Memory-Size    3GB
Dbserver-Disk-Size      10GB
Bootstrapped-At         -
Endpoint-Url            .*
Root-Password           \*\*\* use \'--show-root-password\' to expose \*\*\*
Model                   oneshard
Is-Clone                false
Clone-Backup-Id         
Node-Count              3
Node-Disk-Size          10GB
Node-Size-Id            c4-a4
`
		out, err := tests.RunCommand(args)
		assert.NoError(tt, err)
		assert.True(tt, tests.CompareOutput(out, []byte(compare)))
	})

	t.Run("delete a deployment", func(tt *testing.T) {
		args := []string{"delete", "deployment", "--deployment-id=" + deplId}
		compare := `Deleted deployment!
`
		out, err := tests.RunCommand(args)
		assert.NoError(tt, err)
		assert.True(tt, tests.CompareOutput(out, []byte(compare)))
	})
}
