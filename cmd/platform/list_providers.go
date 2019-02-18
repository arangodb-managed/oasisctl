//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package platform

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	common "github.com/arangodb-managed/apis/common/v1"
	platform "github.com/arangodb-managed/apis/platform/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/format"
)

func init() {
	cmd.InitCommand(
		cmd.ListCmd,
		&cobra.Command{
			Use:   "providers",
			Short: "List all providers the authenticated user has access to",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				cmd.MustCheckNumberOfArgs(args, 0)

				// Connect
				conn := cmd.MustDialAPI()
				platformc := platform.NewPlatformServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Fetch providers
				list, err := platformc.ListProviders(ctx, &common.ListOptions{})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to list providers")
				}

				// Show result
				fmt.Println(format.ProviderList(list.Items, cmd.RootArgs.Format))
			}
		},
	)
}
