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
	// updateRoleCmd updates a role that the user has access to
	updateRoleCmd = &cobra.Command{
		Use:   "role",
		Short: "Update a role the authenticated user has access to",
		Run:   updateRoleCmdRun,
	}
	updateRoleArgs struct {
		roleID            string
		organizationID    string
		name              string
		description       string
		addPermissions    []string
		removePermissions []string
	}
)

func init() {
	cmd.UpdateCmd.AddCommand(updateRoleCmd)
	f := updateRoleCmd.Flags()
	f.StringVarP(&updateRoleArgs.roleID, "role-id", "r", cmd.DefaultRole(), "Identifier of the role")
	f.StringVarP(&updateRoleArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
	f.StringVar(&updateRoleArgs.name, "name", "", "Name of the role")
	f.StringVar(&updateRoleArgs.description, "description", "", "Description of the role")
	f.StringSliceVar(&updateRoleArgs.addPermissions, "add-permission", nil, "Permissions to add to the role")
	f.StringSliceVar(&updateRoleArgs.removePermissions, "remove-permission", nil, "Permissions to remove from the role")
}

func updateRoleCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := updateRoleArgs
	roleID, argsUsed := cmd.OptOption("role-id", cargs.roleID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch role
	item := selection.MustSelectRole(ctx, log, roleID, cargs.organizationID, iamc, rmc)

	// Set changes
	f := c.Flags()
	hasChanges := false
	if f.Changed("name") {
		item.Name = cargs.name
		hasChanges = true
	}
	if f.Changed("description") {
		item.Description = cargs.description
		hasChanges = true
	}
	if len(cargs.addPermissions) > 0 {
		orgLen := len(item.GetPermissions())
		item.Permissions = stringSliceUnion(item.GetPermissions(), cargs.addPermissions)
		hasChanges = hasChanges || len(item.Permissions) != orgLen
	}
	if len(cargs.removePermissions) > 0 {
		orgLen := len(item.GetPermissions())
		item.Permissions = stringSliceExcept(item.GetPermissions(), cargs.removePermissions)
		hasChanges = hasChanges || len(item.Permissions) != orgLen
	}
	if !hasChanges {
		fmt.Println("No changes")
	} else {
		// Update role
		updated, err := iamc.UpdateRole(ctx, item)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to update role")
		}

		// Show result
		fmt.Println("Updated role!")
		fmt.Println(format.Role(updated, cmd.RootArgs.Format))
	}
}

// stringSliceUnion returns a union of the elements in both slices.
func stringSliceUnion(a, b []string) []string {
	m := make(map[string]struct{})
	for _, x := range a {
		m[x] = struct{}{}
	}
	for _, x := range b {
		m[x] = struct{}{}
	}
	result := make([]string, 0, len(m))
	for x := range m {
		result = append(result, x)
	}
	return result
}

// stringSliceExcept returns all elements of a that are not element of b.
func stringSliceExcept(a, b []string) []string {
	m := make(map[string]struct{})
	for _, x := range b {
		m[x] = struct{}{}
	}
	result := make([]string, 0, len(a))
	for x := range m {
		if _, found := m[x]; !found {
			result = append(result, x)
		}
	}
	return result
}
