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
// Author Brautigam Gergely
// Author Ewout Prangsma
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
		cmd.ListCmd,
		&cobra.Command{
			Use:   "examples",
			Short: "List all example datasets",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog

				// Connect
				conn := cmd.MustDialAPI()
				examplec := example.NewExampleDatasetServiceClient(conn)
				ctx := cmd.ContextWithToken()

				list, err := examplec.ListExampleDatasets(ctx, &example.ListExampleDatasetsRequest{})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to list examples")
				}

				// Show result
				fmt.Println(format.ExampleList(list.Items, cmd.RootArgs.Format))
			}
		},
	)
}