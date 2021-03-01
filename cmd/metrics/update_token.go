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

package metrics

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	data "github.com/arangodb-managed/apis/data/v1"
	metrics "github.com/arangodb-managed/apis/metrics/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.UpdateMetricsCmd,
		&cobra.Command{
			Use:   "token",
			Short: "Update a metrics token",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				tokenID        string
				organizationID string
				projectID      string
				deploymentID   string
				name           string
				description    string
			}{}
			f.StringVarP(&cargs.tokenID, "token-id", "t", cmd.DefaultMetricsToken(), "Identifier of the metrics token")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")
			f.StringVarP(&cargs.deploymentID, "deployment-id", "d", cmd.DefaultDeployment(), "Identifier of the deployment")
			f.StringVar(&cargs.name, "name", "", "Name of the CA certificate")
			f.StringVar(&cargs.description, "description", "", "Description of the CA certificate")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				tokenID, argsUsed := cmd.OptOption("token-id", cargs.tokenID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				metricsc := metrics.NewMetricsServiceClient(conn)
				datac := data.NewDataServiceClient(conn)
				rmc := rm.NewResourceManagerServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Fetch token
				item := selection.MustSelectMetricsToken(ctx, log, tokenID, cargs.deploymentID, cargs.projectID, cargs.organizationID, metricsc, datac, rmc)

				// Set changes
				f := c.Flags()
				hasChanges := false
				if f.Changed("name") {
					item.Name = cargs.name
					hasChanges = true
				}
				if f.Changed("description") {
					item.Description = cargs.description
					hasChanges = true
				}
				if !hasChanges {
					fmt.Println("No changes")
				} else {
					// Update token
					updated, err := metricsc.UpdateToken(ctx, item)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to update metrics token")
					}

					// Show result
					fmt.Println("Updated metrics token!")
					fmt.Println(format.MetricsToken(updated, cmd.RootArgs.Format))
				}
			}
		},
	)
}
