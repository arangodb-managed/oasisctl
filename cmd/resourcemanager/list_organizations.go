//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package rm

import (
	"fmt"

	"github.com/spf13/cobra"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/format"
)

var (
	// listOrganizationsCmd fetches organizations the user is a part of
	listOrganizationsCmd = &cobra.Command{
		Use:   "organizations",
		Short: "List all organizations the authenticated user is a member of",
		Run:   listOrganizationsCmdRun,
	}
)

func init() {
	cmd.ListCmd.AddCommand(listOrganizationsCmd)
}

func listOrganizationsCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cmd.MustCheckNumberOfArgs(args, 0)

	// Connect
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch organizations
	list, err := rmc.ListOrganizations(ctx, &common.ListOptions{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list organizations")
	}

	// Show result
	fmt.Println(format.OrganizationList(list.Items, cmd.RootArgs.Format))
}
