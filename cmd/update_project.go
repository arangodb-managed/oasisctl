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
	updateCmd.AddCommand(updateProjectCmd)
	f := updateProjectCmd.Flags()
	f.StringVarP(&updateProjectArgs.projectID, "project-id", "p", defaultProject(), "Identifier of the project")
	f.StringVarP(&updateProjectArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
	f.StringVarP(&updateProjectArgs.name, "name", "n", "", "Name of the project")
	f.StringVarP(&updateProjectArgs.description, "description", "d", "", "Description of the project")
}

func updateProjectCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	projectID, argsUsed := optOption("project-id", updateProjectArgs.projectID, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch project
	item := selection.MustSelectProject(ctx, cliLog, projectID, updateProjectArgs.organizationID, rmc)

	// Set changes
	f := cmd.Flags()
	hasChanges := false
	if f.Changed("name") {
		item.Name = updateProjectArgs.name
		hasChanges = true
	}
	if f.Changed("description") {
		item.Description = updateProjectArgs.description
		hasChanges = true
	}
	if !hasChanges {
		fmt.Println("No changes")
	} else {
		// Update project
		updated, err := rmc.UpdateProject(ctx, item)
		if err != nil {
			cliLog.Fatal().Err(err).Msg("Failed to update project")
		}

		// Show result
		fmt.Println("Updated project!")
		fmt.Println(format.Project(updated, rootArgs.format))
	}
}
