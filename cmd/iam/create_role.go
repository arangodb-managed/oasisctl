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
	// createRoleCmd creates a new role
	createRoleCmd = &cobra.Command{
		Use:   "role",
		Short: "Create a new role",
		Run:   createRoleCmdRun,
	}
	createRoleArgs struct {
		name           string
		description    string
		organizationID string
		permissions    []string
	}
)

func init() {
	cmd.CreateCmd.AddCommand(createRoleCmd)

	f := createRoleCmd.Flags()
	f.StringVarP(&createRoleArgs.name, "name", "n", "", "Name of the role")
	f.StringVarP(&createRoleArgs.description, "description", "d", "", "Description of the role")
	f.StringSliceVarP(&createRoleArgs.permissions, "permission", "p", nil, "Permissions granted by the role")
	f.StringVarP(&createRoleArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization to create the role in")
}

func createRoleCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	name, argsUsed := cmd.ReqOption("name", createRoleArgs.name, args, 0)
	description := createRoleArgs.description
	permissions := createRoleArgs.permissions
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch organization
	org := selection.MustSelectOrganization(ctx, cmd.CLILog, createRoleArgs.organizationID, rmc)

	// Create role
	result, err := iamc.CreateRole(ctx, &iam.Role{
		OrganizationId: org.GetId(),
		Name:           name,
		Description:    description,
		Permissions:    permissions,
	})
	if err != nil {
		cmd.CLILog.Fatal().Err(err).Msg("Failed to create role")
	}

	// Show result
	fmt.Println("Success!")
	fmt.Println(format.Role(result, cmd.RootArgs.Format))
}
