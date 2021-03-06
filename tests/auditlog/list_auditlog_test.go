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

package auditlog

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	audit "github.com/arangodb-managed/apis/audit/v1"
	common "github.com/arangodb-managed/apis/common/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/tests"
)

func TestListAuditLog(t *testing.T) {
	// Initialize the root command.
	cmd.RootCmd.PersistentPreRun(nil, nil)
	ctx := cmd.ContextWithToken()
	conn := cmd.MustDialAPI()
	defer conn.Close()
	auditc := audit.NewAuditServiceClient(conn)
	org, err := tests.GetDefaultOrganization()
	require.NoError(t, err)

	name := "test-auditlog"
	description := "test-description"
	auditLog, err := auditc.CreateAuditLog(ctx, &audit.AuditLog{
		Name:           name,
		Description:    description,
		OrganizationId: org,
	})
	require.NoError(t, err)

	// Cleanup
	defer func() {
		if _, err := auditc.DeleteAuditLog(ctx, &common.IDOptions{Id: auditLog.GetId()}); err != nil {
			t.Log(err)
		}
	}()

	args := []string{"list", "auditlogs", fmt.Sprintf("--organization-id=%s", org)}
	compare := `Id\s+| Name\s+| Url\s+| Description\s+| Default | Created-At\s+| Deleted-At | Destinations(\s.*)*` +
		auditLog.GetId() + `\s+| ` + name + `\s+| /Organization/` + org + `/AuditLog/` + auditLog.GetId() + `\s+| ` + description + `\s+|.*|.*|.*| None
`
	out, err := tests.RunCommand(args)
	require.NoError(t, err)
	assert.True(t, tests.CompareOutput(out, []byte(compare)))
}
