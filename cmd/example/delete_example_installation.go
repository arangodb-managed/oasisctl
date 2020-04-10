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
// Author Ewout Prangsma
//

package example

import (
	"fmt"

	common "github.com/arangodb-managed/apis/common/v1"
	example "github.com/arangodb-managed/apis/example/v1"
	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

func init() {
	cmd.InitCommand(
		DeleteExampleCmd,
		&cobra.Command{
			Use:   "installation",
			Short: "Delete an example datasets installation",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				installationID string
			}{}
			f.StringVar(&cargs.installationID, "installation-id", "", "The ID of the installation to delete.")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				installationID, argsUsed := cmd.ReqOption("installation-id", cargs.installationID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				examplec := example.NewExampleDatasetServiceClient(conn)
				ctx := cmd.ContextWithToken()

				if _, err := examplec.DeleteExampleDatasetInstallation(ctx, &common.IDOptions{Id: installationID}); err != nil {
					log.Fatal().Err(err).Msg("Failed to delete examples")
				}

				// Show result
				fmt.Println("Success")
			}
		},
	)
}
