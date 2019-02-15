//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package iam

import (
	"fmt"

	"github.com/spf13/cobra"

	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/selection"
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
	cmd.DeleteCmd.AddCommand(deleteGroupCmd)
	f := deleteGroupCmd.Flags()
	f.StringVarP(&deleteGroupArgs.groupID, "group-id", "g", cmd.DefaultGroup(), "Identifier of the group")
	f.StringVarP(&deleteGroupArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
}

func deleteGroupCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	groupID, argsUsed := cmd.OptOption("group-id", deleteGroupArgs.groupID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch group
	item := selection.MustSelectGroup(ctx, cmd.CLILog, groupID, deleteGroupArgs.organizationID, iamc, rmc)

	// Delete group
	if _, err := iamc.DeleteGroup(ctx, &common.IDOptions{Id: item.GetId()}); err != nil {
		cmd.CLILog.Fatal().Err(err).Msg("Failed to delete group")
	}

	// Show result
	fmt.Println("Deleted group!")
}
