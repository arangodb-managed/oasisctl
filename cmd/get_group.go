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
	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/pkg/format"
)

var (
	// getGroupCmd fetches a group that the user has access to
	getGroupCmd = &cobra.Command{
		Use:   "group",
		Short: "Get a group the authenticated user has access to",
		Run:   getGroupCmdRun,
	}
	getGroupArgs struct {
		groupID        string
		organizationID string
	}
)

func init() {
	getCmd.AddCommand(getGroupCmd)
	f := getGroupCmd.Flags()
	f.StringVarP(&getGroupArgs.groupID, "group-id", "g", defaultGroup(), "Identifier of the group")
	f.StringVarP(&getGroupArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func getGroupCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	groupID, argsUsed := optOption("group-id", getGroupArgs.groupID, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch group
	item := mustSelectGroup(ctx, groupID, getGroupArgs.organizationID, iamc, rmc)

	// Show result
	fmt.Println(format.Group(item, rootArgs.format))
}

// mustSelectGroup fetches the group with given ID.
// If no ID is specified, all groups are fetched from the selected organization
// and if the list is exactly 1 long, that group is returned.
func mustSelectGroup(ctx context.Context, id, orgID string, iamc iam.IAMServiceClient, rmc rm.ResourceManagerServiceClient) *iam.Group {
	if id == "" {
		org := mustSelectOrganization(ctx, orgID, rmc)
		list, err := iamc.ListGroups(ctx, &common.ListOptions{ContextId: org.GetId()})
		if err != nil {
			cliLog.Fatal().Err(err).Msg("Failed to list groups")
		}
		if len(list.Items) != 1 {
			cliLog.Fatal().Err(err).Msg("You have access to %d groups. Please specify one explicitly.")
		}
		return list.Items[0]
	}
	result, err := iamc.GetGroup(ctx, &common.IDOptions{Id: id})
	if err != nil {
		cliLog.Fatal().Err(err).Str("group", id).Msg("Failed to get group")
	}
	return result
}
