//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/pkg/format"
)

var (
	// getOrganizationCmd fetches an organization the user is a part of
	getOrganizationCmd = &cobra.Command{
		Use:   "organization",
		Short: "Get an organization the authenticated user is a member of",
		Run:   getOrganizationCmdRun,
	}
	getOrganizationArgs struct {
		organizationID string
	}
)

func init() {
	getCmd.AddCommand(getOrganizationCmd)
	f := getOrganizationCmd.Flags()
	f.StringVarP(&getProjectsArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func getOrganizationCmdRun(cmd *cobra.Command, args []string) {
	// Connect
	conn := mustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch organization
	item := mustSelectOrganization(ctx, getOrganizationArgs.organizationID, rmc)

	// Show result
	fmt.Println(format.Organization(item, rootArgs.format))
}

// mustSelectOrganization fetches the organization with given ID.
// If no ID is specified, all organizations are fetched and if the user
// is member of exactly 1, that organization is returned.
func mustSelectOrganization(ctx context.Context, id string, rmc rm.ResourceManagerServiceClient) *rm.Organization {
	if id == "" {
		list, err := rmc.ListOrganizations(ctx, &common.ListOptions{})
		if err != nil {
			cliLog.Fatal().Err(err).Msg("Failed to list organizations")
		}
		if len(list.Items) != 1 {
			cliLog.Fatal().Err(err).Msg("You're member of %d organizations. Please specify one explicitly.")
		}
		return list.Items[0]
	}
	result, err := rmc.GetOrganization(ctx, &common.IDOptions{Id: id})
	if err != nil {
		cliLog.Fatal().Err(err).Str("organization", id).Msg("Failed to get organization")
	}
	return result
}
