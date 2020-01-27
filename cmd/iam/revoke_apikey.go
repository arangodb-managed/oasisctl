//
// DISCLAIMER
//
// Copyright 2019 ArangoDB GmbH, Cologne, Germany
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
