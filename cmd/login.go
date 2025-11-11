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

package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	iam "github.com/arangodb-managed/apis/iam/v1"
)

var (
	// LoginCmd is used to login using an API key
	LoginCmd = &cobra.Command{
		Use: "login",
		Run: ShowUsage,
	}
)

func init() {
	InitCommand(
		RootCmd,
		&cobra.Command{
			Use:   "login",
			Short: "Log in to the Arango Managed Platform (AMP) using an API key",
			Long: `To authenticate in a script environment, run:
	
	export OASIS_TOKEN=$(oasisctl login --key-id=<your-key-id> --key-secret=<your-key-secret>)
`,
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				keyID     string
				keySecret string
			}{}
			f.StringVarP(&cargs.keyID, "key-id", "i", "", "API key identifier")
			f.StringVarP(&cargs.keySecret, "key-secret", "s", "", "API key secret")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := CLILog
				keyID, argsUsed := ReqOption("key-id", cargs.keyID, args, 0)
				keySecret, argsUsed := ReqOption("key-secret", cargs.keySecret, args, argsUsed)
				MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := MustDialAPI()
				iamc := iam.NewIAMServiceClient(conn)
				ctx := context.Background()

				resp, err := iamc.AuthenticateAPIKey(ctx, &iam.AuthenticateAPIKeyRequest{
					Id:     keyID,
					Secret: keySecret,
				})
				if err != nil {
					log.Fatal().Err(err).Msg("Authentication failed")
				}
				fmt.Println(resp.GetToken())
			}
		},
	)
}
