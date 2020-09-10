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

func TestCreateDeployment(t *testing.T) {
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
	args := []string{"create", "deployment", "--name=" + deplName, "--organization-id=" + org, "--project-id=" + proj, "--region-id=" + region, "--version=" + version.GetVersion(), "--accept"}

	defer func() {
		list, err := datac.ListDeployments(ctx, &common.ListOptions{ContextId: proj})
		if err != nil {
			t.Log(err)
			return
		}
		for _, d := range list.GetItems() {
			if _, err := datac.DeleteDeployment(ctx, &common.IDOptions{Id: d.GetId()}); err != nil {
				t.Log(err)
			}
		}
	}()

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
}
