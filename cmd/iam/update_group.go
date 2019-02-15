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

	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/format"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

var (
	// updateGroupCmd updates a group that the user has access to
	updateGroupCmd = &cobra.Command{
		Use:   "group",
		Short: "Update a group the authenticated user has access to",
		Run:   updateGroupCmdRun,
	}
	updateGroupArgs struct {
		groupID        string
		organizationID string
		name           string
		description    string
	}
)

func init() {
	cmd.UpdateCmd.AddCommand(updateGroupCmd)
	f := updateGroupCmd.Flags()
	f.StringVarP(&updateGroupArgs.groupID, "group-id", "g", cmd.DefaultGroup(), "Identifier of the group")
	f.StringVarP(&updateGroupArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
	f.StringVarP(&updateGroupArgs.name, "name", "n", "", "Name of the group")
	f.StringVarP(&updateGroupArgs.description, "description", "d", "", "Description of the group")
}

func updateGroupCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	groupID, argsUsed := cmd.OptOption("group-id", updateGroupArgs.groupID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch group
	item := selection.MustSelectGroup(ctx, cmd.CLILog, groupID, updateGroupArgs.organizationID, iamc, rmc)

	// Set changes
	f := c.Flags()
	hasChanges := false
	if f.Changed("name") {
		item.Name = updateGroupArgs.name
		hasChanges = true
	}
	if f.Changed("description") {
		item.Description = updateGroupArgs.description
		hasChanges = true
	}
	if !hasChanges {
		fmt.Println("No changes")
	} else {
		// Update group
		updated, err := iamc.UpdateGroup(ctx, item)
		if err != nil {
			cmd.CLILog.Fatal().Err(err).Msg("Failed to update group")
		}

		// Show result
		fmt.Println("Updated group!")
		fmt.Println(format.Group(updated, cmd.RootArgs.Format))
	}
}
