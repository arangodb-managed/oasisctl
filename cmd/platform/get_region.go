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

	platform "github.com/arangodb-managed/apis/platform/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/format"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.GetCmd,
		&cobra.Command{
			Use:   "region",
			Short: "Get a region the authenticated user has access to",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				organizationID string
				regionID   string
				providerID string
			}{}
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Optional Identifier of the organization")
			f.StringVarP(&cargs.regionID, "region-id", "r", cmd.DefaultRegion(), "Identifier of the region")
			f.StringVarP(&cargs.providerID, "provider-id", "p", cmd.DefaultProvider(), "Identifier of the provider")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				regionID, argsUsed := cmd.OptOption("region-id", cargs.regionID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				platformc := platform.NewPlatformServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Fetch region
				item := selection.MustSelectRegion(ctx, log, regionID, cargs.providerID, cargs.organizationID, platformc)

				// Show result
				fmt.Println(format.Region(item, cmd.RootArgs.Format))
			}
		},
	)
}
