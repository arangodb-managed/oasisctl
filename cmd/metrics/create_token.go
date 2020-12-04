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

package metrics

import (
	"fmt"

	"github.com/dchest/uniuri"
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
		cmd.CreateMetricsCmd,
		&cobra.Command{
			Use:   "token",
			Short: "Create a new metrics access token",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				name           string
				description    string
				organizationID string
				projectID      string
				deploymentID   string
			}{}
			f.StringVar(&cargs.name, "name", "", "Name of the token")
			f.StringVar(&cargs.description, "description", "", "Description of the token")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization to create the token in")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project to create the token in")
			f.StringVarP(&cargs.deploymentID, "deployment-id", "d", cmd.DefaultDeployment(), "Identifier of the deployment to create the token for")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				deploymentID, argsUsed := cmd.OptOption("deployment-id", cargs.deploymentID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)
				name := cargs.name
				if name == "" {
					name = "token-" + uniuri.NewLen(4)
				}

				// Connect
				conn := cmd.MustDialAPI()
				metricsc := metrics.NewMetricsServiceClient(conn)
				datac := data.NewDataServiceClient(conn)
				rmc := rm.NewResourceManagerServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Fetch deployment
				deployment := selection.MustSelectDeployment(ctx, log, deploymentID, cargs.projectID, cargs.organizationID, datac, rmc)

				// Create token
				result, err := metricsc.CreateToken(ctx, &metrics.Token{
					DeploymentId: deployment.GetId(),
					Name:         name,
					Description:  cargs.description,
				})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to create metrics token")
				}

				// Show result
				format.DisplaySuccess(cmd.RootArgs.Format)
				fmt.Println(format.MetricsToken(result, cmd.RootArgs.Format))
			}
		},
	)
}
