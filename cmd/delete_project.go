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
)

var (
	// deleteProjectCmd deletes a project that the user has access to
	deleteProjectCmd = &cobra.Command{
		Use:   "project",
		Short: "Delete a project the authenticated user has access to",
		Run:   deleteProjectCmdRun,
	}
	deleteProjectArgs struct {
		organizationID string
		projectID      string
	}
)

func init() {
	deleteCmd.AddCommand(deleteProjectCmd)
	f := deleteProjectCmd.Flags()
	f.StringVarP(&deleteProjectArgs.projectID, "project-id", "p", defaultProject(), "Identifier of the project")
	f.StringVarP(&deleteProjectArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func deleteProjectCmdRun(cmd *cobra.Command, args []string) {
	// Connect
	conn := mustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch project
	item := mustSelectProject(ctx, deleteProjectArgs.projectID, deleteProjectArgs.organizationID, rmc)

	// Delete project
	if _, err := rmc.DeleteProject(ctx, item); err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to delete project")
	}

	// Show result
	fmt.Println("Deleted project!")
}
