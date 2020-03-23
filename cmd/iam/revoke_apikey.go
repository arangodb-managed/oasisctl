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

package iam

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
)

func init() {
	cmd.InitCommand(
		cmd.RevokeCmd,
		&cobra.Command{
			Use:   "apikey",
			Short: "Revoke an API key with given identifier",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				apiKeyID string
			}{}
			f.StringVarP(&cargs.apiKeyID, "apikey-id", "i", "", "Identifier of the API key to revoke")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				apiKeyID, argsUsed := cmd.ReqOption("apikey-id", cargs.apiKeyID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				iamc := iam.NewIAMServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Revoke API key
				_, err := iamc.RevokeAPIKey(ctx, &common.IDOptions{
					Id: apiKeyID,
				})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to revoke API key")
				}

				// Show result
				fmt.Println("Revoked API key!")
			}
		},
	)
}
