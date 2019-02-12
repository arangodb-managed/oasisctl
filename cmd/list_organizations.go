//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

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
	listCmd.AddCommand(listOrganizationsCmd)
}

func listOrganizationsCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	mustCheckNumberOfArgs(args, 0)

	// Connect
	conn := mustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch organizations
	list, err := rmc.ListOrganizations(ctx, &common.ListOptions{})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to list organizations")
	}

	// Show result
	fmt.Println(format.OrganizationList(list.Items, rootArgs.format))
}
