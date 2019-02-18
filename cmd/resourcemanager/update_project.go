//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package rm

import (
	"fmt"
	"github.com/arangodb-managed/oasis/cmd"

	"github.com/spf13/cobra"

	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
	"github.com/arangodb-managed/oasis/pkg/format"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

var (
	// updateProjectCmd updates a project that the user has access to
	updateProjectCmd = &cobra.Command{
		Use:   "project",
		Short: "Update a project the authenticated user has access to",
		Run:   updateProjectCmdRun,
	}
	updateProjectArgs struct {
		projectID      string
		organizationID string
		name           string
		description    string
	}
)

func init() {
	cmd.UpdateCmd.AddCommand(updateProjectCmd)
	f := updateProjectCmd.Flags()
	f.StringVarP(&updateProjectArgs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")
	f.StringVarP(&updateProjectArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
	f.StringVar(&updateProjectArgs.name, "name", "", "Name of the project")
	f.StringVar(&updateProjectArgs.description, "description", "", "Description of the project")
}

func updateProjectCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := updateProjectArgs
	projectID, argsUsed := cmd.OptOption("project-id", cargs.projectID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch project
	item := selection.MustSelectProject(ctx, log, projectID, cargs.organizationID, rmc)

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
	if !hasChanges {
		fmt.Println("No changes")
	} else {
		// Update project
		updated, err := rmc.UpdateProject(ctx, item)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to update project")
		}

		// Show result
		fmt.Println("Updated project!")
		fmt.Println(format.Project(updated, cmd.RootArgs.Format))
	}
}
