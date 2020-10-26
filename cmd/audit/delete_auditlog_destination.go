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
	common "github.com/arangodb-managed/apis/common/v1"

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
				auditlogType   string
				url            string
				organizationID string
			}{}
			f.StringVarP(&cargs.id, "auditlog-id", "i", "", "Identifier of the auditlog to delete.")
			f.IntVar(&cargs.index, "index", -1, "Index of the destination to remove. Only one delete option can be specified.")
			f.StringVar(&cargs.auditlogType, "type", "", "Type of the destination to remove. This will remove ALL destinations with that type.")
			f.StringVar(&cargs.url, "url", "", "An optional URL which will be used to delete a single destination instead of all.")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization.")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				id, argsUsed := cmd.ReqOption("auditlog-id", cargs.id, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				auditc := audit.NewAuditServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Make the call
				item := selection.MustSelectAuditLog(ctx, log, id, cargs.organizationID, auditc)
				destinations := item.GetDestinations()

				var err error
				if cargs.index != -1 {
					destinations, err = deleteByIndex(destinations, cargs.index)
					if err != nil {
						log.Fatal().Int("index", cargs.index).Int("length", len(destinations)).Err(err).Msg("Failed to delete destination.")
					}
				} else if cargs.auditlogType != "" {
					destinations, err = deleteByType(destinations, cargs.auditlogType, cargs.url)
					if err != nil {
						log.Fatal().Err(err).Str("type", cargs.auditlogType).Str("url", cargs.url).Msg("Failed to delete destination.")
					}
				}
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

func deleteByIndex(destinations []*audit.AuditLog_Destination, index int) ([]*audit.AuditLog_Destination, error) {
	if index >= len(destinations) {
		return nil, common.InvalidArgument("The index is larger than the length of destinations.")
	}
	destinations = append(destinations[:index], destinations[index+1:]...)
	return destinations, nil
}

func deleteByType(destinations []*audit.AuditLog_Destination, dType string, url string) ([]*audit.AuditLog_Destination, error) {
	for i := 0; i < len(destinations); i++ {
		if destinations[i].GetType() == dType {
			if url != "" {
				if destinations[i].GetHttpPost().GetUrl() == url {
					destinations = append(destinations[:i], destinations[i+1:]...)
					return destinations, nil
				}
			} else {
				destinations = append(destinations[:i], destinations[i+1:]...)
				i--
			}
		}
	}
	if url != "" {
		return nil, common.NotFound("Destination with given URL not found.")
	}
	return destinations, nil
}
