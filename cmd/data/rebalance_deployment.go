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
// Author Tomasz Mielech
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
		cmd.RebalanceDeploymentCmd,
		&cobra.Command{
			Use:   "shards",
			Short: "Rebalance shards of a deployment",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				deploymentID   string
			}{}
			f.StringVarP(&cargs.deploymentID, "deployment-id", "d", cmd.DefaultDeployment(),
				"Identifier of the deployment")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				deploymentID, argsUsed := cmd.OptOption("deployment-id", cargs.deploymentID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI(cmd.WithoutVersionCheck())
				datac := data.NewDataServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Request rebalance shards
				_, err := datac.RebalanceDeploymentShards(ctx, &data.RebalanceDeploymentShardsRequest{
					DeploymentId: deploymentID,
				})
				if err != nil {
					log.Fatal().Msg("Rebalancing shards requests failed.")
				}

				fmt.Printf("Rebalancing shards has been requested for the deployment '%s'.\n", deploymentID)
			}
		},
	)
}

