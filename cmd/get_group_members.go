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
	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/pkg/format"
)

var (
	// getGroupMembersCmd fetches the members of a group the user is a part of
	getGroupMembersCmd = &cobra.Command{
		Use:   "members",
		Short: "Get members of a group the authenticated user is a member of",
		Run:   getGroupMembersCmdRun,
	}
	getGroupMembersArgs struct {
		groupID        string
		organizationID string
	}
)

func init() {
	getGroupCmd.AddCommand(getGroupMembersCmd)
	f := getGroupMembersCmd.Flags()
	f.StringVarP(&getGroupMembersArgs.groupID, "group-id", "g", defaultGroup(), "Identifier of the group")
	f.StringVarP(&getGroupMembersArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func getGroupMembersCmdRun(cmd *cobra.Command, args []string) {
	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch group
	group := mustSelectGroup(ctx, getGroupMembersArgs.groupID, getGroupMembersArgs.organizationID, iamc, rmc)

	list, err := iamc.ListGroupMembers(ctx, &common.ListOptions{ContextId: group.GetId()})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to list group members")
	}

	// Show result
	fmt.Println(format.GroupMemberList(ctx, list.GetItems(), iamc, rootArgs.format))
}
