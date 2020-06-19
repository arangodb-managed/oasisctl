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

	iam "github.com/arangodb-managed/apis/iam/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
)

func init() {
	cmd.InitCommand(
		revokeAPIKeyCmd,
		&cobra.Command{
			Use:   "token",
			Short: "Revoke an API key token",
			Long:  "Revoke the token (resulting from API key authentication)",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				cmd.MustCheckNumberOfArgs(args, 0)

				// Connect
				conn := cmd.MustDialAPI()
				iamc := iam.NewIAMServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Revoke API key token
				_, err := iamc.RevokeAPIKeyToken(ctx, &iam.RevokeAPIKeyTokenRequest{
					Token: cmd.RootArgs.Token,
				})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to revoke API key token")
				}

				// Show result
				fmt.Println("Revoked API key token!")
			}
		},
	)
}
