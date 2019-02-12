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
	createCmd.AddCommand(createRoleCmd)

	f := createRoleCmd.Flags()
	f.StringVarP(&createRoleArgs.name, "name", "n", "", "Name of the role")
	f.StringVarP(&createRoleArgs.description, "description", "d", "", "Description of the role")
	f.StringSliceVarP(&createRoleArgs.permissions, "permission", "p", nil, "Permissions granted by the role")
	f.StringVarP(&createRoleArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization to create the role in")
}

func createRoleCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	name, argsUsed := reqOption("name", createRoleArgs.name, args, 0)
	description := createRoleArgs.description
	permissions := createRoleArgs.permissions
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch organization
	org := mustSelectOrganization(ctx, createRoleArgs.organizationID, rmc)

	// Create role
	result, err := iamc.CreateRole(ctx, &iam.Role{
		OrganizationId: org.GetId(),
		Name:           name,
		Description:    description,
		Permissions:    permissions,
	})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to create role")
	}

	// Show result
	fmt.Println("Success!")
	fmt.Println(format.Role(result, rootArgs.format))
}
