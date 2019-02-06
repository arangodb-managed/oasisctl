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
	conn := mustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()
	list, err := rmc.ListOrganizations(ctx, &common.ListOptions{})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to list organizations")
	}
	fmt.Printf("%v\n", list)
}
