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
//

package example

import (
	"fmt"

	"github.com/spf13/cobra"

	example "github.com/arangodb-managed/apis/example/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
)

var (
	// listExampleDatasetCmd fetches example datasets
	listExampleDatasetCmd = &cobra.Command{
		Use:   "examples",
		Short: "List all example datasets",
		Run:   listExampleDatasetCmdRun,
	}
)

func init() {
	cmd.ListCmd.AddCommand(listExampleDatasetCmd)
}

func listExampleDatasetCmdRun(c *cobra.Command, args []string) {
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
