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
)

var (
	// deleteGroupCmd deletes a group that the user has access to
	deleteGroupCmd = &cobra.Command{
		Use:   "group",
		Short: "Delete a group the authenticated user has access to",
		Run:   deleteGroupCmdRun,
	}
	deleteGroupArgs struct {
		organizationID string
		groupID        string
	}
)

func init() {
	deleteCmd.AddCommand(deleteGroupCmd)
	f := deleteGroupCmd.Flags()
	f.StringVarP(&deleteGroupArgs.groupID, "group-id", "g", defaultGroup(), "Identifier of the group")
	f.StringVarP(&deleteGroupArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func deleteGroupCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	groupID, argsUsed := optOption("group-id", deleteGroupArgs.groupID, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch group
	item := mustSelectGroup(ctx, groupID, deleteGroupArgs.organizationID, iamc, rmc)

	// Delete group
	if _, err := iamc.DeleteGroup(ctx, &common.IDOptions{Id: item.GetId()}); err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to delete group")
	}

	// Show result
	fmt.Println("Deleted group!")
}
