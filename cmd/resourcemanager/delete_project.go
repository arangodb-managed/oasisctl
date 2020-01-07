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

	"github.com/spf13/cobra"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
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
	cmd.DeleteCmd.AddCommand(deleteProjectCmd)
	f := deleteProjectCmd.Flags()
	f.StringVarP(&deleteProjectArgs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")
	f.StringVarP(&deleteProjectArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
}

func deleteProjectCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := deleteProjectArgs
	projectID, argsUsed := cmd.OptOption("project-id", cargs.projectID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch project
	item := selection.MustSelectProject(ctx, log, projectID, cargs.organizationID, rmc)

	// Delete project
	if _, err := rmc.DeleteProject(ctx, &common.IDOptions{Id: item.GetId()}); err != nil {
		log.Fatal().Err(err).Msg("Failed to delete project")
	}

	// Show result
	fmt.Println("Deleted project!")
}
