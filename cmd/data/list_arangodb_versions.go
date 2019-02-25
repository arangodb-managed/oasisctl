//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package data

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	common "github.com/arangodb-managed/apis/common/v1"
	data "github.com/arangodb-managed/apis/data/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/format"
)

func init() {
	cmd.InitCommand(
		cmd.ListArangoDBCmd,
		&cobra.Command{
			Use:   "versions",
			Short: "List all supported ArangoDB versions",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			//cargs := &struct {}{}

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				cmd.MustCheckNumberOfArgs(args, 0)

				// Connect
				conn := cmd.MustDialAPI()
				datac := data.NewDataServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Fetch versions
				list, err := datac.ListVersions(ctx, &common.ListOptions{})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to list versions")
				}

				// Fetch default version
				defaultVersion, err := datac.GetDefaultVersion(ctx, &common.Empty{})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to get default version")
				}

				// Show result
				fmt.Println(format.VersionList(list.Items, defaultVersion, cmd.RootArgs.Format))
			}
		},
	)
}
