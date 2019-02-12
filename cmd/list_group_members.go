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
	// listGroupMembersCmd fetches the members of a group the user is a part of
	listGroupMembersCmd = &cobra.Command{
		Use:   "members",
		Short: "List members of a group the authenticated user is a member of",
		Run:   listGroupMembersCmdRun,
	}
	listGroupMembersArgs struct {
		groupID        string
		organizationID string
	}
)

func init() {
	listGroupCmd.AddCommand(listGroupMembersCmd)
	f := listGroupMembersCmd.Flags()
	f.StringVarP(&listGroupMembersArgs.groupID, "group-id", "g", defaultGroup(), "Identifier of the group")
	f.StringVarP(&listGroupMembersArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func listGroupMembersCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	groupID, argsUsed := reqOption("group-id", listGroupMembersArgs.groupID, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch group
	group := mustSelectGroup(ctx, groupID, listGroupMembersArgs.organizationID, iamc, rmc)

	list, err := iamc.ListGroupMembers(ctx, &common.ListOptions{ContextId: group.GetId()})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to list group members")
	}

	// Show result
	fmt.Println(format.GroupMemberList(ctx, list.GetItems(), iamc, rootArgs.format))
}
