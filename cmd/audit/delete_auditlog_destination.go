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
		deleteAuditLogCmd,
		&cobra.Command{
			Use:   "destination",
			Short: "Delete a destination from an auditlog",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				id             string
				name           string
				index          int
				organizationID string
			}{}
			f.StringVarP(&cargs.id, "auditlog-id", "i", "", "Identifier of the auditlog to delete.")
			f.IntVar(&cargs.index, "index", -1, "Index of the destination to remove.")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				id, argsUsed := cmd.ReqOption("auditlog-id", cargs.id, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				if cargs.index < 0 {
					log.Fatal().Msg("Please provide a valid index.")
				}

				// Connect
				conn := cmd.MustDialAPI()
				auditc := audit.NewAuditServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Make the call
				item := selection.MustSelectAuditLog(ctx, log, id, cargs.organizationID, auditc)
				destinations := item.GetDestinations()
				if cargs.index >= len(destinations) {
					log.Fatal().Msg("The index is larger than the length of destinations.")
				}
				destinations = append(destinations[:cargs.index], destinations[cargs.index+1:]...)
				item.Destinations = destinations

				// Update auditlog with the new destinations.
				result, err := auditc.UpdateAuditLog(ctx, item)
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to delete destination.")
				}

				// Show result
				format.DisplaySuccess(cmd.RootArgs.Format)
				fmt.Println(format.AuditLog(result, cmd.RootArgs.Format))
			}
		},
	)
}
