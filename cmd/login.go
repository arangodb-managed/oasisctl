//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
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
			Short: "Login to ArangoDB Oasis using an API key",
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
