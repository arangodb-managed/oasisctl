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

	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/cmd"
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
	cmd.GetCmd.AddCommand(getProjectCmd)
	f := getProjectCmd.Flags()
	f.StringVarP(&getProjectArgs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")
	f.StringVarP(&getProjectArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
}

func getProjectCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	projectID, argsUsed := cmd.OptOption("project-id", getProjectArgs.projectID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch project
	item := selection.MustSelectProject(ctx, cmd.CLILog, projectID, getProjectArgs.organizationID, rmc)

	// Show result
	fmt.Println(format.Project(item, cmd.RootArgs.Format))
}