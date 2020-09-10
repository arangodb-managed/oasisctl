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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/tests"
)

func TestListProject(t *testing.T) {
	cmd.RootCmd.PersistentPreRun(nil, nil)
	ctx := cmd.ContextWithToken()
	conn := cmd.MustDialAPI()
	org, err := tests.GetDefaultOrganization()
	require.NoError(t, err)
	rmc := rm.NewResourceManagerServiceClient(conn)

	project, err := rmc.CreateProject(ctx, &rm.Project{Name: "testProject", OrganizationId: org})
	assert.NoError(t, err)
	defer func() {
		if _, err := rmc.DeleteProject(ctx, &common.IDOptions{Id: project.GetId()}); err != nil {
			t.Log(err)
		}
	}()

	// We must acount for the default organization
	// (\s.*\s.*)+? -> match any number of rows in between lazily
	compare := `^Id        | Name\s+| Description\s+| Url\s+| Created-At(\s.*)*` + project.GetId() + ` | testProject\s+|.*| ` + project.GetUrl() + ` | .*`
	args := []string{"list", "projects", "--organization-id=" + org}
	out, err := tests.RunCommand(args)
	require.NoError(t, err)
	assert.True(t, tests.CompareOutput(out, []byte(compare)))
}
