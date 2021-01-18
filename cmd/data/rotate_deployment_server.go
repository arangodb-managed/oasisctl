//
// DISCLAIMER
//
// Copyright 2021 ArangoDB GmbH, Cologne, Germany
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

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	data "github.com/arangodb-managed/apis/data/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
)

func init() {
	cmd.InitCommand(
		cmd.RotateDeploymentCmd,
		&cobra.Command{
			Use:   "server",
			Short: "Rotate a single server of a deployment",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				serverIDs      []string
				deploymentID   string
				organizationID string
				projectID      string
			}{}
			f.StringSliceVarP(&cargs.serverIDs, "server-id", "s", cmd.SplitByComma(cmd.DefaultServer()), "Identifier of the deployment server")
			f.StringVarP(&cargs.deploymentID, "deployment-id", "d", cmd.DefaultDeployment(), "Identifier of the deployment")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				deploymentID, argsUsed := cmd.OptOption("deployment-id", cargs.deploymentID, args, 0)
				serverIDs, argsUsed := cmd.OptOptionSlice("server-id", cargs.serverIDs, args, argsUsed)
				cmd.MustCheckNumberOfArgs(args, argsUsed)
				if len(serverIDs) == 0 {
					log.Fatal().Msg("Missing server ID(s)")
				}

				// Connect
				conn := cmd.MustDialAPI()
				datac := data.NewDataServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Request server rotation
				errors := 0
				for _, serverID := range serverIDs {
					if _, err := datac.RotateDeploymentServer(ctx, &data.RotateDeploymentServerRequest{
						DeploymentId: deploymentID,
						ServerId:     serverID,
					}); err != nil {
						errors++
						log.Error().Err(err).Msg("Failed to rotate server.")
					} else {
						fmt.Printf("Server rotation has been requested for server '%s'.\n", serverID)
					}
				}
				if errors > 0 {
					log.Fatal().Msg("One or more server rotation requests failed.")
				}
			}
		},
	)
}
