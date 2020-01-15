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

	iam "github.com/arangodb-managed/apis/iam/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
)

var (
	// deleteGroupMembersCmd deletes a list of members from a group
	deleteGroupMembersCmd = &cobra.Command{
		Use:   "group",
		Short: "Add members to group",
		Run:   deleteGroupMembersCmdRun,
	}
	deleteGroupMembersArgs struct {
		groupID string
		userIDs []string
	}
)

func init() {
	cmd.DeleteCmd.AddCommand(deleteGroupMembersCmd)

	f := deleteGroupMembersCmd.Flags()
	f.StringVarP(&deleteGroupMembersArgs.groupID, "group-id", "g", cmd.DefaultGroup(), "Identifier of the group")
	deleteGroupMembersArgs.userIDs = *f.StringSliceP("user-ids", "u", []string{}, "A coma separated list of user ids")

}

func deleteGroupMembersCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := deleteGroupMembersArgs
	groupID, argsUsed := cmd.OptOption("group-id", cargs.groupID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Delete members
	_, err := iamc.DeleteGroupMembers(ctx, &iam.GroupMembersRequest{
		GroupId: groupID,
		UserIds: cargs.userIDs,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to delete users.")
	}

	fmt.Println("Success!")
}
