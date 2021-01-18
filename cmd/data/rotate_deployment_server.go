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
				serverID       string
				deploymentID   string
				organizationID string
				projectID      string
			}{}
			f.StringVarP(&cargs.serverID, "server-id", "s", cmd.DefaultServer(), "Identifier of the deployment server")
			f.StringVarP(&cargs.deploymentID, "deployment-id", "d", cmd.DefaultDeployment(), "Identifier of the deployment")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				deploymentID, argsUsed := cmd.OptOption("deployment-id", cargs.deploymentID, args, 0)
				serverID, argsUsed := cmd.OptOption("deployment-id", cargs.deploymentID, args, argsUsed)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				datac := data.NewDataServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Request server rotation
				if _, err := datac.RotateDeploymentServer(ctx, &data.RotateDeploymentServerRequest{
					DeploymentId: deploymentID,
					ServerId:     serverID,
				}); err != nil {
					log.Fatal().Err(err).Msg("Failed to rotate server.")
				}

				fmt.Println("Server rotation has been requested.")
			}
		},
	)
}
