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

package cmd

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/arangodb-managed/oasisctl/pkg/expectedapis"
)

const (
	apisJSONFilename = "apis.json"
)

func init() {
	InitCommand(
		RootCmd,
		&cobra.Command{
			Use:    "expected-apis",
			Short:  "Generate an apis.json file.",
			Long:   "Generates a file which contains all the versions needed by this tool.",
			Hidden: true,
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			c.Run = func(c *cobra.Command, args []string) {
				log := CLILog
				content, err := json.MarshalIndent(expectedapis.ExpectedAPIVersions(), "", "  ")
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to marshal map to json.")
				}
				if err := os.WriteFile(apisJSONFilename, content, os.ModePerm); err != nil {
					log.Fatal().Err(err).Msg("Failed to write out file.")
				}
			}
		},
	)
}
