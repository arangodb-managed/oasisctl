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

	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
	"github.com/arangodb-managed/oasis/pkg/format"
)

var (
	// createProjectCmd creates a new project
	createProjectCmd = &cobra.Command{
		Use:   "project",
		Short: "Create a new project",
		Run:   createProjectCmdRun,
	}
	createProjectArgs struct {
		name           string
		description    string
		organizationID string
	}
)

func init() {
	createCmd.AddCommand(createProjectCmd)

	f := createProjectCmd.Flags()
	f.StringVarP(&createProjectArgs.name, "name", "n", "", "Name of the project")
	f.StringVarP(&createProjectArgs.description, "description", "d", "", "Description of the project")
	f.StringVarP(&createProjectArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization to create the project in")
}

func createProjectCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	name, argsUsed := reqOption("name", createProjectArgs.name, args, 0)
	description := createProjectArgs.description
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch organization
	org := mustSelectOrganization(ctx, createProjectArgs.organizationID, rmc)

	// Create project
	result, err := rmc.CreateProject(ctx, &rm.Project{
		OrganizationId: org.GetId(),
		Name:           name,
		Description:    description,
	})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to create project")
	}

	// Show result
	fmt.Println("Success!")
	fmt.Println(format.Project(result, rootArgs.format))
}
