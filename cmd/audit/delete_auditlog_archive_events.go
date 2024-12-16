//
// DISCLAIMER
//
// Copyright 2020-2024 ArangoDB GmbH, Cologne, Germany
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

package audit

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/timestamppb"

	audit "github.com/arangodb-managed/apis/audit/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
	"github.com/arangodb-managed/oasisctl/pkg/util"
)

func init() {
	cmd.InitCommand(
		deleteAuditArchive,
		&cobra.Command{
			Use:   "events",
			Short: "Delete auditlog archive events",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				auditLogArchiveID string
				to                string
			}{}
			f.StringVarP(&cargs.auditLogArchiveID, "auditlog-archive-id", "i", "", "Identifier of the auditlog archive to delete events from.")
			f.StringVar(&cargs.to, "to", "", "Remove events created before this timestamp.")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				id, argsUsed := cmd.ReqOption("auditlog-archive-id", cargs.auditLogArchiveID, args, 0)
				to, argsUsed := cmd.OptOption("to", cargs.to, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				auditc := audit.NewAuditServiceClient(conn)
				ctx := cmd.ContextWithToken()

				var (
					toDate *timestamppb.Timestamp
					err    error
				)
				if to != "" {
					toDate, err = util.ParseTime(to)
					if err != nil {
						log.Fatal().Err(err).Str("date", to).Msg("Failed to parse to timestamp.")
					}
				}
				auditLogArchive := selection.MustSelectAuditLogArchive(ctx, log, id, auditc)
				// Make the call
				if _, err := auditc.DeleteAuditLogArchiveEvents(ctx, &audit.DeleteAuditLogArchiveEventsRequest{
					AuditlogarchiveId: auditLogArchive.GetId(),
					To:                toDate,
				}); err != nil {
					log.Fatal().Err(err).Msg("Failed to delete audit archive events log.")
				}

				// Show result
				fmt.Println("Deleted auditlog archive events!")
			}
		},
	)
}
