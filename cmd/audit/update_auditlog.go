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

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	audit "github.com/arangodb-managed/apis/audit/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.UpdateCmd,
		&cobra.Command{
			Use:   "auditlog",
			Short: "Update an auditlog",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				auditlogID     string
				name           string
				description    string
				isDefault      bool
				organizationID string
			}{}
			f.StringVarP(&cargs.auditlogID, "auditlog-id", "i", "", "Identifier of the auditlog to update.")
			f.StringVar(&cargs.name, "name", "", "Name of the audit log.")
			f.StringVar(&cargs.description, "description", "", "Description of the audit log.")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
			f.BoolVar(&cargs.isDefault, "default", false, "If set, this AuditLog is the default for the organization.")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				auditLogID, argsUsed := cmd.ReqOption("auditlog-id", cargs.auditlogID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				auditc := audit.NewAuditServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Construct the default destination
				item := selection.MustSelectAuditLog(ctx, log, auditLogID, cargs.organizationID, auditc)

				f := c.Flags()
				hasChanges := false
				changeDefault := false
				newDefaultAuditLogID := ""
				if f.Changed("name") {
					item.Name = cargs.name
					hasChanges = true
				}
				if f.Changed("description") {
					item.Description = cargs.description
					hasChanges = true
				}
				if f.Changed("default") {
					if item.GetIsDefault() != cargs.isDefault {
						changeDefault = true
						hasChanges = true
						if cargs.isDefault {
							newDefaultAuditLogID = item.GetId()
						} else {
							newDefaultAuditLogID = ""
						}
					}
				}
				if !hasChanges {
					fmt.Println("No changes")
					return
				}

				// Update audit log
				result, err := auditc.UpdateAuditLog(ctx, item)
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to update auditlog.")
				}

				// Update default audit log (if requested)
				if changeDefault {
					if _, err := auditc.SetDefaultAuditLog(ctx, &audit.SetDefaultAuditLogRequest{
						OrganizationId: item.GetOrganizationId(),
						AuditlogId:     newDefaultAuditLogID,
					}); err != nil {
						log.Fatal().Err(err).Msg("Failed to change default audit log.")
					}
				}

				// Show result
				format.DisplaySuccess(cmd.RootArgs.Format)
				fmt.Println(format.AuditLog(result, cmd.RootArgs.Format))
			}
		},
	)
}
