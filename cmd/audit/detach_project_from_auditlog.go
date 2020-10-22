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
	audit "github.com/arangodb-managed/apis/audit/v1"
	common "github.com/arangodb-managed/apis/common/v1"
	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

func init() {
	cmd.InitCommand(
		cmd.AuditLogCmd,
		&cobra.Command{
			Use:   "detach",
			Short: "Detach a project from an auditlog",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				projectID string
			}{}
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project to detach the audit log.")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				projID, argsUsed := cmd.ReqOption("project-id", cargs.projectID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				auditc := audit.NewAuditServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Make the call
				if _, err := auditc.DetachProjectFromAuditLog(ctx, &common.IDOptions{Id: projID}); err != nil {
					log.Fatal().Err(err).Str("project-id", projID).Msg("Failed to detach to project.")
				}

				// Show result
				format.DisplaySuccess(cmd.RootArgs.Format)
			}
		},
	)
}
