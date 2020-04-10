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

package example

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	example "github.com/arangodb-managed/apis/example/v1"
	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
)

func init() {
	cmd.InitCommand(
		CreateExampleCmd,
		&cobra.Command{
			Use:   "installation",
			Short: "Create a new example dataset installation",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				deploymentID     string
				exampleDatasetID string
				description      string
			}{}
			f.StringVarP(&cargs.deploymentID, "deployment-id", "d", cmd.DefaultDeployment(), "Identifier of the deployment to list installations for")
			f.StringVar(&cargs.exampleDatasetID, "example-dataset-id", "", "ID of the example dataset")
			f.StringVar(&cargs.description, "description", "", "Description of the installation")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				deploymentID, argsUsed := cmd.ReqOption("deployment-id", cargs.deploymentID, args, 0)
				exampleDatasetID, argsUsed := cmd.ReqOption("example-dataset-id", cargs.exampleDatasetID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				examplec := example.NewExampleDatasetServiceClient(conn)
				ctx := cmd.ContextWithToken()

				req := &example.ExampleDatasetInstallation{
					DeploymentId:     deploymentID,
					ExampledatasetId: exampleDatasetID,
				}
				result, err := examplec.CreateExampleDatasetInstallation(ctx, req)

				if err != nil {
					log.Fatal().Err(err).Msg("Failed to create installation")
				}

				// Show result
				fmt.Println("Success!")
				fmt.Println(format.ExampleDatasetInstallation(result, cmd.RootArgs.Format))
			}
		},
	)
}
