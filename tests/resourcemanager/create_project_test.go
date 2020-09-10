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

package resourcemanager

import (
	"testing"

	"github.com/arangodb-managed/oasisctl/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
)

func TestCreateProject(t *testing.T) {
	cmd.RootCmd.PersistentPreRun(nil, nil)
	ctx := cmd.ContextWithToken()
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	org, err := tests.GetDefaultOrganization()
	require.NoError(t, err)

	testProj := "testCreateProject"
	defer func() {
		list, err := rmc.ListProjects(ctx, &common.ListOptions{ContextId: org})
		if err != nil {
			t.Log(err)
		}

		// We only delete the test organization
		for _, project := range list.GetItems() {
			if project.GetName() == testProj {
				if _, err := rmc.DeleteProject(ctx, &common.IDOptions{Id: project.GetId()}); err != nil {
					t.Log(err)
				}
				break
			}
		}
	}()

	compare := `^Success!
Id          \d+
Name        ` + testProj + `
Description 
Url         /Organization/\d+/Project/.*
Created-At  .*
Deleted-At  -
$`

	args := []string{"create", "project", "--name=" + testProj, "--organization-id=" + org}
	out, err := tests.RunCommand(args)
	require.NoError(t, err)
	assert.True(t, tests.CompareOutput(out, []byte(compare)))
}
