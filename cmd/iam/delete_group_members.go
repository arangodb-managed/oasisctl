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
		Short: "Delete members from group",
		Run:   deleteGroupMembersCmdRun,
	}
	deleteGroupMembersArgs struct {
		groupID        string
		organizationID string
		userEmails     *[]string
	}
)

func init() {
	deleteGroupCmd.AddCommand(deleteGroupMembersCmd)

	f := deleteGroupMembersCmd.Flags()
	f.StringVarP(&deleteGroupMembersArgs.groupID, "group-id", "g", cmd.DefaultGroup(), "Identifier of the group to delete members from")
	f.StringVarP(&deleteGroupMembersArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
	deleteGroupMembersArgs.userEmails = f.StringSliceP("user-emails", "u", []string{}, "A comma separated list of user email addresses")
}

func deleteGroupMembersCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := deleteGroupMembersArgs
	groupID, argsUsed := cmd.OptOption("group-id", cargs.groupID, args, 0)
	organizationID, argsUsed := cmd.OptOption("organization-id", cargs.organizationID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	ctx := cmd.ContextWithToken()
	rmc := rm.NewResourceManagerServiceClient(conn)

	organization := selection.MustSelectGroup(ctx, log, groupID, organizationID, iamc, rmc)
	group := selection.MustSelectGroup(ctx, log, groupID, organization.Id, iamc, rmc)

	log.Info().Msgf("Deleting members: %s", cargs.userEmails)
	var userIds []string
	members, err := iamc.ListGroupMembers(ctx, &common.ListOptions{ContextId: group.Id})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list group members.")
	}
	emailIDMap := make(map[string]string)
	for _, id := range members.Items {
		user, err := iamc.GetUser(ctx, &common.IDOptions{Id: id})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get user")
		}
		emailIDMap[user.Email] = user.Id
	}

	for _, e := range *cargs.userEmails {
		if id, ok := emailIDMap[e]; !ok {
			log.Fatal().Str("email", e).Str("group-id", group.Id).Msg("User not part of the group")
		} else {
			userIds = append(userIds, id)
		}
	}

	if _, err := iamc.DeleteGroupMembers(ctx, &iam.GroupMembersRequest{GroupId: group.Id, UserIds: userIds}); err != nil {
		log.Fatal().Err(err).Msg("Failed to delete users.")
	}

	fmt.Println("Success!")
}
