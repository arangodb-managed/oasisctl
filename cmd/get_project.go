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
	// getProjectCmd fetches a project that the user has access to
	getProjectCmd = &cobra.Command{
		Use:   "project",
		Short: "Get a project the authenticated user has access to",
		Run:   getProjectCmdRun,
	}
	getProjectArgs struct {
		organizationID string
		projectID      string
	}
)

func init() {
	getCmd.AddCommand(getProjectCmd)
	f := getProjectCmd.Flags()
	f.StringVarP(&getProjectArgs.projectID, "project-id", "p", defaultProject(), "Identifier of the project")
	f.StringVarP(&getProjectArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func getProjectCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	projectID, argsUsed := optOption("project-id", getProjectArgs.projectID, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch project
	item := selection.MustSelectProject(ctx, cliLog, projectID, getProjectArgs.organizationID, rmc)

	// Show result
	fmt.Println(format.Project(item, rootArgs.format))
}
