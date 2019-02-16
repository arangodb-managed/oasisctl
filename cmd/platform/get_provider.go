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
	// getProviderCmd fetches a cloud provider that the user has access to
	getProviderCmd = &cobra.Command{
		Use:   "provider",
		Short: "Get a provider the authenticated user has access to",
		Run:   getProviderCmdRun,
	}
	getProviderArgs struct {
		providerID string
	}
)

func init() {
	cmd.GetCmd.AddCommand(getProviderCmd)
	f := getProviderCmd.Flags()
	f.StringVarP(&getProviderArgs.providerID, "provider-id", "p", cmd.DefaultProvider(), "Identifier of the provider")
}

func getProviderCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := getProviderArgs
	providerID, argsUsed := cmd.OptOption("provider-id", cargs.providerID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	platformc := platform.NewPlatformServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch provider
	item := selection.MustSelectProvider(ctx, log, providerID, platformc)

	// Show result
	fmt.Println(format.Provider(item, cmd.RootArgs.Format))
}
