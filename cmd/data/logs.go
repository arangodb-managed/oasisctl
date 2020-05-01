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
// Author Ewout Prangsma
//

package data

import (
	"fmt"
	"io"

	types "github.com/gogo/protobuf/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	data "github.com/arangodb-managed/apis/data/v1"
	mon "github.com/arangodb-managed/apis/monitoring/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
	"github.com/arangodb-managed/oasisctl/pkg/util"
)

func init() {
	cmd.InitCommand(
		cmd.RootCmd,
		&cobra.Command{
			Use:   "logs",
			Short: "Get logs of the servers of a deployment the authenticated user has access to",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				deploymentID   string
				organizationID string
				projectID      string
				role           string
				limit          int
				start          string
				end            string
			}{}
			f.StringVarP(&cargs.deploymentID, "deployment-id", "d", cmd.DefaultDeployment(), "Identifier of the deployment")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")
			f.StringVarP(&cargs.role, "role", "r", "", "Limit logs to servers with given role only (agents|coordinators|dbservers)")
			f.IntVarP(&cargs.limit, "limit", "l", 0, "Limit the number of log lines")
			f.StringVar(&cargs.start, "start", "", "Start fetching logs from this timestamp (pass timestamp or duration before now)")
			f.StringVar(&cargs.end, "end", "", "End fetching logs at this timestamp (pass timestamp or duration before now)")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				deploymentID, argsUsed := cmd.OptOption("deployment-id", cargs.deploymentID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				datac := data.NewDataServiceClient(conn)
				rmc := rm.NewResourceManagerServiceClient(conn)
				monc := mon.NewMonitoringServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Fetch deployment
				item := selection.MustSelectDeployment(ctx, log, deploymentID, cargs.projectID, cargs.organizationID, datac, rmc)

				// Fetch logs
				req := &mon.GetDeploymentLogsRequest{
					DeploymentId: item.GetId(),
					Role:         cargs.role,
					Limit:        int32(cargs.limit),
				}
				if ts := cargs.start; ts != "" {
					t, err := util.ParseTimeFromNow(ts)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to parse start time")
					}
					req.StartAt, err = types.TimestampProto(t)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to encode start time")
					}
				}
				if ts := cargs.end; ts != "" {
					t, err := util.ParseTimeFromNow(ts)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to parse end time")
					}
					req.EndAt, err = types.TimestampProto(t)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to encode end time")
					}
				}
				client, err := monc.GetDeploymentLogs(ctx, req)
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to fetch deployment logs")
				}
				log.Debug().Msg("GetDeploymentLogs succeeded")

				// Show logs
				for {
					msg, err := client.Recv()
					if err == io.EOF {
						// All done
						break
					} else if err != nil {
						log.Fatal().Err(err).Msg("Failed to next deployment logs chunk")
					}
					fmt.Print(string(msg.GetChunk()))
				}
			}
		},
	)
}
