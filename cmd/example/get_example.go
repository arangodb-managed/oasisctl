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

package example

import (
	"fmt"

	example "github.com/arangodb-managed/apis/example/v1"
	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

var getExampleCmd = cmd.InitCommand(
	cmd.GetCmd,
	&cobra.Command{
		Use:   "example",
		Short: "Get a single example dataset",
	},
	func(c *cobra.Command, f *flag.FlagSet) {
		cargs := &struct {
			exampleDatasetID string
		}{}
		f.StringVarP(&cargs.exampleDatasetID, "example-dataset-id", "e", "", "ID of the example dataset")

		c.Run = func(c *cobra.Command, args []string) {
			// Validate arguments
			log := cmd.CLILog
			exampleDatasetID, argsUsed := cmd.OptOption("example-dataset-id", cargs.exampleDatasetID, args, 0)
			cmd.MustCheckNumberOfArgs(args, argsUsed)

			// Connect
			conn := cmd.MustDialAPI()
			examplec := example.NewExampleDatasetServiceClient(conn)
			ctx := cmd.ContextWithToken()

			// Select example
			example := selection.MustSelectExampleDataset(ctx, log, exampleDatasetID, examplec)

			// Show result
			fmt.Println(format.Example(example, cmd.RootArgs.Format))
		}
	},
)
