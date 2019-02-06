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

	"github.com/arangodb-managed/adbcloud/pkg/format"
)

var (
	// getOrganizationsCmd fetches organizations the user is a part of
	getOrganizationsCmd = &cobra.Command{
		Use:   "organizations",
		Short: "Get all organizations the authenticated user is a member of",
		Run:   getOrganizationsCmdRun,
	}
)

func init() {
	getCmd.AddCommand(getOrganizationsCmd)
}

func getOrganizationsCmdRun(cmd *cobra.Command, args []string) {
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
