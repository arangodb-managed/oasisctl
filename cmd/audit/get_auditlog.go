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

package audit

import (
	"fmt"

	"github.com/arangodb-managed/oasisctl/pkg/selection"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	audit "github.com/arangodb-managed/apis/audit/v1"
	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
)

var getAuditLogCmd = cmd.InitCommand(
	cmd.GetCmd,
	&cobra.Command{
		Use:   "auditlog",
		Short: "Get auditlog archive",
	},
	func(c *cobra.Command, f *flag.FlagSet) {
		cargs := &struct {
			auditLogID     string
			organizationID string
		}{}
		f.StringVarP(&cargs.auditLogID, "auditlog-id", "i", "", "Identifier of the auditlog")
		f.StringVarP(&cargs.auditLogID, "organization-id", "o", "", "Identifier of the organization")

		c.Run = func(c *cobra.Command, args []string) {
			// Validate arguments
			log := cmd.CLILog
			auditLogId, argsUsed := cmd.ReqOption("auditlog-id", cargs.auditLogID, args, 0)
			cmd.MustCheckNumberOfArgs(args, argsUsed)

			// Connect
			conn := cmd.MustDialAPI()
			auditc := audit.NewAuditServiceClient(conn)
			ctx := cmd.ContextWithToken()
			item := selection.MustSelectAuditLog(ctx, log, auditLogId, cargs.organizationID, auditc)

			// Show result
			format.DisplaySuccess(cmd.RootArgs.Format)
			fmt.Println(format.AuditLog(item, cmd.RootArgs.Format))
		}
	},
)
