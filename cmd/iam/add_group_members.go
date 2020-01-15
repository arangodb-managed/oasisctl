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
	// addGroupMembersCmd adds a list of members to a group
	addGroupMembersCmd = &cobra.Command{
		Use:   "group",
		Short: "Add members to group",
		Run:   addGroupMembersCmdRun,
	}
	addGroupMembersArgs struct {
		groupID string
		userIDs []string
	}
)

func init() {
	cmd.CreateCmd.AddCommand(addGroupMembersCmd)

	f := addGroupMembersCmd.Flags()
	f.StringVarP(&addGroupMembersArgs.groupID, "group-id", "g", cmd.DefaultGroup(), "Identifier of the group")
	addGroupMembersArgs.userIDs = *f.StringSliceP("user-ids", "u", []string{}, "A coma separated list of user ids")

}

func addGroupMembersCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	// Validate arguments
	log := cmd.CLILog
	cargs := addGroupMembersArgs
	groupID, argsUsed := cmd.OptOption("group-id", cargs.groupID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Add members
	_, err := iamc.AddGroupMembers(ctx, &iam.GroupMembersRequest{
		GroupId: groupID,
		UserIds: cargs.userIDs,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to add users.")
	}

	fmt.Println("Success!")
}
