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

	common "github.com/arangodb-managed/apis/common/v1"
	platform "github.com/arangodb-managed/apis/platform/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/format"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

var (
	// listRegionsCmd fetches region of the given provider
	listRegionsCmd = &cobra.Command{
		Use:   "regions",
		Short: "List all regions of the given provider",
		Run:   listRegionsCmdRun,
	}
	listRegionsArgs struct {
		providerID string
	}
)

func init() {
	cmd.ListCmd.AddCommand(listRegionsCmd)
	f := listRegionsCmd.Flags()
	f.StringVarP(&listRegionsArgs.providerID, "provider-id", "p", cmd.DefaultProvider(), "Identifier of the provider")
}

func listRegionsCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := listRegionsArgs
	providerID, argsUsed := cmd.OptOption("provider-id", cargs.providerID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	platformc := platform.NewPlatformServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch provider
	provider := selection.MustSelectProvider(ctx, log, providerID, platformc)

	// Fetch regions in provider
	list, err := platformc.ListRegions(ctx, &common.ListOptions{ContextId: provider.GetId()})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list regions")
	}

	// Show result
	fmt.Println(format.RegionList(list.Items, cmd.RootArgs.Format))
}
