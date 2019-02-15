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

	platform "github.com/arangodb-managed/apis/platform/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/format"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

var (
	// getRegionCmd fetches a region that the user has access to
	getRegionCmd = &cobra.Command{
		Use:   "region",
		Short: "Get a region the authenticated user has access to",
		Run:   getRegionCmdRun,
	}
	getRegionArgs struct {
		regionID   string
		providerID string
	}
)

func init() {
	cmd.GetCmd.AddCommand(getRegionCmd)
	f := getRegionCmd.Flags()
	f.StringVarP(&getRegionArgs.regionID, "region-id", "r", cmd.DefaultRegion(), "Identifier of the region")
	f.StringVarP(&getRegionArgs.providerID, "provider-id", "p", cmd.DefaultProvider(), "Identifier of the provider")
}

func getRegionCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	regionID, argsUsed := cmd.OptOption("region-id", getRegionArgs.regionID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	platformc := platform.NewPlatformServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch region
	item := selection.MustSelectRegion(ctx, cmd.CLILog, regionID, getRegionArgs.providerID, platformc)

	// Show result
	fmt.Println(format.Region(item, cmd.RootArgs.Format))
}
