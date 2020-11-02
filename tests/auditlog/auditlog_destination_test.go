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

func TestAddAuditLogDestination(t *testing.T) {
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
		Destinations: []*audit.AuditLog_Destination{
			{
				Type: audit.DestinationCloud,
			},
		},
	})
	require.NoError(t, err)
	defer func() {
		if _, err := auditc.DeleteAuditLog(ctx, &common.IDOptions{Id: auditLog.GetId()}); err != nil {
			t.Log(err)
		}
	}()

	args := []string{"add", "auditlog", "destination", fmt.Sprintf("--auditlog-id=%s", auditLog.GetId()), fmt.Sprintf("--destination-type=%s", audit.DestinationCloud)}
	compare := `Success!
Id           ` + auditLog.GetId() + `
Name         ` + auditLog.GetName() + `
Url          /Organization/` + org + `/AuditLog/` + auditLog.GetId() + `
Description  ` + auditLog.GetDescription() + `
Default      -
Created-At   .*
Deleted-At   -
Destinations Index | Type
0     | cloud
1     | cloud
`
	out, err := tests.RunCommand(args)
	require.NoError(t, err)
	assert.True(t, tests.CompareOutput(out, []byte(compare)))

	args = []string{"delete", "auditlog", "destination", "--auditlog-id=" + auditLog.GetId(), "--index=1"}
	compare = `Success!
Id           ` + auditLog.GetId() + `
Name         ` + auditLog.GetName() + `
Url          /Organization/` + org + `/AuditLog/` + auditLog.GetId() + `
Description  ` + auditLog.GetDescription() + `
Default      -
Created-At   .*
Deleted-At   -
Destinations Index | Type
0     | cloud
`

	out, err = tests.RunCommand(args)
	require.NoError(t, err)
	assert.True(t, tests.CompareOutput(out, []byte(compare)))

	args = []string{"list", "auditlog", "destinations", "--auditlog-id=" + auditLog.GetId()}
	compare = `Index | Type
0     | cloud
`

	out, err = tests.RunCommand(args)
	require.NoError(t, err)
	assert.True(t, tests.CompareOutput(out, []byte(compare)))

	auditLog2, err := auditc.CreateAuditLog(ctx, &audit.AuditLog{
		Name:           name,
		Description:    description,
		OrganizationId: org,
		Destinations: []*audit.AuditLog_Destination{
			{
				Type: audit.DestinationCloud,
			},
			{
				Type: audit.DestinationHTTPSPost,
				HttpPost: &audit.AuditLog_HttpsPostSettings{
					Url: "https://google.com/asdf1",
				},
			},
			{
				Type: audit.DestinationHTTPSPost,
				HttpPost: &audit.AuditLog_HttpsPostSettings{
					Url: "https://google.com/asdf2",
				},
			},
			{
				Type: audit.DestinationHTTPSPost,
				HttpPost: &audit.AuditLog_HttpsPostSettings{
					Url: "https://google.com/asdf3",
				},
			},
		},
	})
	require.NoError(t, err)
	defer func() {
		if _, err := auditc.DeleteAuditLog(ctx, &common.IDOptions{Id: auditLog2.GetId()}); err != nil {
			t.Log(err)
		}
	}()
	args = []string{"delete", "auditlog", "destination", "--auditlog-id=" + auditLog2.GetId(), "--type=https-post", "--url=https://google.com/asdf2"}
	compare = `Success!
Id           ` + auditLog2.GetId() + `
Name         ` + auditLog2.GetName() + `
Url          /Organization/` + org + `/AuditLog/` + auditLog2.GetId() + `
Description  ` + auditLog2.GetDescription() + `
Default      -
Created-At   .*
Deleted-At   -
Destinations Index | Type
0     | cloud
1     | https-post | https://google.com/asdf1 | 
2     | https-post | https://google.com/asdf3 |
`

	out, err = tests.RunCommand(args)
	require.NoError(t, err)
	assert.True(t, tests.CompareOutput(out, []byte(compare)))

	args = []string{"delete", "auditlog", "destination", "--auditlog-id=" + auditLog2.GetId(), "--type=https-post"}
	compare = `Success!
Id           ` + auditLog2.GetId() + `
Name         ` + auditLog2.GetName() + `
Url          /Organization/` + org + `/AuditLog/` + auditLog2.GetId() + `
Description  ` + auditLog2.GetDescription() + `
Default      -
Created-At   .*
Deleted-At   -
Destinations Index | Type
0     | cloud
`

	out, err = tests.RunCommand(args)
	require.NoError(t, err)
	assert.True(t, tests.CompareOutput(out, []byte(compare)))
}
