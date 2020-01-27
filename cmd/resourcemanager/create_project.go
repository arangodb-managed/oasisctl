//
// DISCLAIMER
//
// Copyright 2019 ArangoDB GmbH, Cologne, Germany
//
// Author Ewout Prangsma
//

package rm

import (
	"fmt"

	"github.com/spf13/cobra"

	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
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
	cmd.CreateCmd.AddCommand(createProjectCmd)

	f := createProjectCmd.Flags()
	f.StringVar(&createProjectArgs.name, "name", "", "Name of the project")
	f.StringVar(&createProjectArgs.description, "description", "", "Description of the project")
	f.StringVarP(&createProjectArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization to create the project in")
}

func createProjectCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := createProjectArgs
	name, argsUsed := cmd.ReqOption("name", cargs.name, args, 0)
	description := createProjectArgs.description
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch organization
	org := selection.MustSelectOrganization(ctx, log, cargs.organizationID, rmc)

	// Create project
	result, err := rmc.CreateProject(ctx, &rm.Project{
		OrganizationId: org.GetId(),
		Name:           name,
		Description:    description,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create project")
	}

	// Show result
	fmt.Println("Success!")
	fmt.Println(format.Project(result, cmd.RootArgs.Format))
}
