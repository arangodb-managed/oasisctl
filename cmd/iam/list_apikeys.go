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
	"github.com/arangodb-managed/oasisctl/pkg/format"
)

func init() {
	cmd.InitCommand(
		cmd.ListCmd,
		&cobra.Command{
			Use:   "apikeys",
			Short: "List all API keys created for the current user",
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

				// List API keys
				result, err := iamc.ListAPIKeys(ctx, &common.ListOptions{})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to list API keys")
				}

				// Show result
				fmt.Println(format.APIKeyList(result.GetItems(), cmd.RootArgs.Format))
			}
		},
	)
}
