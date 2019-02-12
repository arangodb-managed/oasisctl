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

	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
	"github.com/arangodb-managed/oasis/pkg/format"
)

var (
	// updateRoleCmd updates a role that the user has access to
	updateRoleCmd = &cobra.Command{
		Use:   "role",
		Short: "Update a role the authenticated user has access to",
		Run:   updateRoleCmdRun,
	}
	updateRoleArgs struct {
		roleID         string
		organizationID string
		name           string
		description    string
	}
)

func init() {
	updateCmd.AddCommand(updateRoleCmd)
	f := updateRoleCmd.Flags()
	f.StringVarP(&updateRoleArgs.roleID, "role-id", "r", defaultRole(), "Identifier of the role")
	f.StringVarP(&updateRoleArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
	f.StringVarP(&updateRoleArgs.name, "name", "n", "", "Name of the role")
	f.StringVarP(&updateRoleArgs.description, "description", "d", "", "Description of the role")
}

func updateRoleCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	roleID, argsUsed := optOption("role-id", updateRoleArgs.roleID, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch role
	item := mustSelectRole(ctx, roleID, updateRoleArgs.organizationID, iamc, rmc)

	// Set changes
	f := cmd.Flags()
	hasChanges := false
	if f.Changed("name") {
		item.Name = updateRoleArgs.name
		hasChanges = true
	}
	if f.Changed("description") {
		item.Description = updateRoleArgs.description
		hasChanges = true
	}
	if !hasChanges {
		fmt.Println("No changes")
	} else {
		// Update role
		updated, err := iamc.UpdateRole(ctx, item)
		if err != nil {
			cliLog.Fatal().Err(err).Msg("Failed to update role")
		}

		// Show result
		fmt.Println("Updated role!")
		fmt.Println(format.Role(updated, rootArgs.format))
	}
}
