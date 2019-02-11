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
	// createGroupCmd creates a new group
	createGroupCmd = &cobra.Command{
		Use:   "group",
		Short: "Create a new group",
		Run:   createGroupCmdRun,
	}
	createGroupArgs struct {
		name           string
		description    string
		organizationID string
	}
)

func init() {
	createCmd.AddCommand(createGroupCmd)

	f := createGroupCmd.Flags()
	f.StringVarP(&createGroupArgs.name, "name", "n", "", "Name of the group")
	f.StringVarP(&createGroupArgs.description, "description", "d", "", "Description of the group")
	f.StringVarP(&createGroupArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization to create the group in")
}

func createGroupCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	name := reqOption("name", createGroupArgs.name, args, 0)
	description := createGroupArgs.description

	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch organization
	org := mustSelectOrganization(ctx, createGroupArgs.organizationID, rmc)

	// Create group
	result, err := iamc.CreateGroup(ctx, &iam.Group{
		OrganizationId: org.GetId(),
		Name:           name,
		Description:    description,
	})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to create group")
	}

	// Show result
	fmt.Println("Success!")
	fmt.Println(format.Group(result, rootArgs.format))
}
