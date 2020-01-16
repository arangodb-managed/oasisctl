//
// DISCLAIMER
//
// Copyright 2020 ArangoDB Inc, Cologne, Germany
//
// Author Gergely Brautigam
//

package iam

import (
	"fmt"

	"github.com/spf13/cobra"

	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

var (
	// deleteGroupMembersCmd deletes a list of members from a group
	deleteGroupMembersCmd = &cobra.Command{
		Use:   "members",
		Short: "Add members to group",
		Run:   deleteGroupMembersCmdRun,
	}
	deleteGroupMembersArgs struct {
		organizationID string
		groupID        string
		userEmails     []string
	}
)

func init() {
	deleteMembersCmd.AddCommand(deleteGroupMembersCmd)

	f := deleteGroupMembersCmd.Flags()
	f.StringVarP(&deleteGroupMembersArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
	f.StringVarP(&deleteGroupMembersArgs.groupID, "group-id", "g", cmd.DefaultGroup(), "Identifier of the group to delete members from")
	deleteGroupMembersArgs.userEmails = *f.StringSliceP("user-emails", "u", []string{}, "A comma separated list of user email addresses")
}

func deleteGroupMembersCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := deleteGroupMembersArgs
	groupID, argsUsed := cmd.OptOption("group-id", cargs.groupID, args, 0)
	organizationID, argsUsed := cmd.OptOption("organization-id", cargs.groupID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	ctx := cmd.ContextWithToken()
	rmc := rm.NewResourceManagerServiceClient(conn)

	var userIds []string
	members, err := rmc.ListOrganizationMembers(ctx, &common.ListOptions{ContextId: organizationID})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list organization members.")
	}
	emailIDMap, err := selection.GenerateUserEmailMap(ctx, members, iamc)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to find user.")
	}

	for _, e := range cargs.userEmails {
		if id, ok := emailIDMap[e]; !ok {
			log.Fatal().Str("email", e).Msg("User not found or not part of the ogranization")
		} else {
			if resp, err := iamc.IsMemberOfGroup(ctx, &iam.IsMemberOfGroupRequest{UserId: id, GroupId: groupID}); err != nil {
				log.Fatal().Err(err).Str("email", e).Str("group-id", groupID).Msgf("Failed to determine if user is member of the group.")
			} else if !resp.Result {
				log.Fatal().Err(err).Str("email", e).Str("group-id", groupID).Msgf("User is not a member of group.")
			}
			userIds = append(userIds, id)
		}
	}

	if _, err := iamc.DeleteGroupMembers(ctx, &iam.GroupMembersRequest{GroupId: groupID, UserIds: userIds}); err != nil {
		log.Fatal().Err(err).Msg("Failed to delete users.")
	}

	fmt.Println("Success!")
}
