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
	"io"

	"github.com/gogo/protobuf/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	audit "github.com/arangodb-managed/apis/audit/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/util"
)

func init() {
	cmd.InitCommand(
		cmd.GetAuditLogCmd,
		&cobra.Command{
			Use:   "events",
			Short: "Get auditlog events",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				auditLogID        string
				limit             int
				auditLogArchiveID string
				from              string
				to                string
				includedTopics    []string
				excludedTopics    []string
			}{}
			f.StringVarP(&cargs.auditLogID, "auditlog-id", "i", "", "Identifier of the auditlog")
			f.StringVar(&cargs.auditLogArchiveID, "auditlog-archive-id", "", "If set, include only events from this AuditLogArchive")
			f.IntVar(&cargs.limit, "limit", 0, "Limit the number of audit log events. Defaults to 0, meaning no limit")
			f.StringVar(&cargs.from, "from", "", "Request events created at or after this timestamp")
			f.StringVar(&cargs.to, "to", "", "Request events created before this timestamp")
			f.StringSliceVar(&cargs.includedTopics, "included-topics", nil, "If non-empty, only request events with one of these topics")
			f.StringSliceVar(&cargs.excludedTopics, "excluded-topics", nil, "If non-empty, leave out events with one of these topics. This takes priority over included")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				auditLogId, argsUsed := cmd.ReqOption("auditlog-id", cargs.auditLogID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				auditc := audit.NewAuditServiceClient(conn)
				ctx := cmd.ContextWithToken()

				var (
					toDate   *types.Timestamp
					fromDate *types.Timestamp
					err      error
				)

				if cargs.to != "" {
					toDate, err = util.ParseTime(cargs.to)
					if err != nil {
						log.Fatal().Err(err).Str("date", cargs.to).Msg("Failed to parse to date.")
					}
				}

				if cargs.from != "" {
					fromDate, err = util.ParseTime(cargs.from)
					if err != nil {
						log.Fatal().Err(err).Str("date", cargs.from).Msg("Failed to parse from date.")
					}
				}

				// Make the call
				result, err := auditc.GetAuditLogEvents(ctx, &audit.GetAuditLogEventsRequest{
					AuditlogId:        auditLogId,
					AuditlogarchiveId: cargs.auditLogArchiveID,
					From:              fromDate,
					To:                toDate,
					IncludedTopics:    cargs.includedTopics,
					ExcludedTopics:    cargs.excludedTopics,
					Limit:             int32(cargs.limit),
				})
				if err != nil {
					log.Fatal().Err(err).Str("auditlog-id", auditLogId).Msg("Failed to get auditlog events.")
				}
				// We don't display success as the stream could break in between.
				for {
					events, err := result.Recv()
					if err == io.EOF {
						// Connection closed normally, retry connection
						break
					}
					if err != nil {
						log.Fatal().Err(err).Msg("Error while loading new events.")
					}

					fmt.Println(format.AuditLogEventList(events.GetItems(), cmd.RootArgs.Format))
				}
			}
		},
	)
}
