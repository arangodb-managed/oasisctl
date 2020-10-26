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

func TestDeleteAuditLog(t *testing.T) {
	org, err := tests.GetDefaultOrganization()
	require.NoError(t, err)
	conn := cmd.MustDialAPI()
	defer conn.Close()
	ctx := cmd.ContextWithToken()
	auditc := audit.NewAuditServiceClient(conn)
	name := "test-auditlog"
	description := "test-description"
	auditLog, err := auditc.CreateAuditLog(ctx, &audit.AuditLog{
		Name:           name,
		Description:    description,
		OrganizationId: org,
	})
	require.NoError(t, err)

	args := []string{"delete", "auditlog", fmt.Sprintf("-i=%s", auditLog.GetId())}
	compare := `^Deleted AuditLog!
$`

	out, err := tests.RunCommand(args)
	require.NoError(t, err)
	assert.True(t, tests.CompareOutput(out, []byte(compare)))

	_, err = auditc.GetAuditLog(ctx, &common.IDOptions{Id: auditLog.GetId()})
	require.Error(t, err)
	assert.True(t, common.IsNotFound(err))
}
