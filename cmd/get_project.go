//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/pkg/format"
)

var (
	// getProjectCmd fetches a project that the user has access to
	getProjectCmd = &cobra.Command{
		Use:   "project",
		Short: "Get a project the authenticated user has access to",
		Run:   getProjectCmdRun,
	}
	getProjectArgs struct {
		projectID string
	}
)

func init() {
	getCmd.AddCommand(getProjectCmd)
	f := getProjectCmd.Flags()
	f.StringVarP(&getProjectArgs.projectID, "project-id", "p", defaultProject(), "Identifier of the project")
}

func getProjectCmdRun(cmd *cobra.Command, args []string) {
	// Connect
	conn := mustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch project
	item := mustSelectProject(ctx, getProjectArgs.projectID, rmc)

	// Show result
	fmt.Println(format.Project(item, rootArgs.format))
}

// mustSelectProject fetches the project with given ID.
// If no ID is specified, all projects are fetched from the selected organization
// and if the list is exactly 1 long, that project is returned.
func mustSelectProject(ctx context.Context, id string, rmc rm.ResourceManagerServiceClient) *rm.Project {
	if id == "" {
		org := mustSelectOrganization(ctx, "", rmc)
		list, err := rmc.ListProjects(ctx, &common.ListOptions{ContextId: org.GetId()})
		if err != nil {
			cliLog.Fatal().Err(err).Msg("Failed to list projects")
		}
		if len(list.Items) != 1 {
			cliLog.Fatal().Err(err).Msg("You have access to %d projects. Please specify one explicitly.")
		}
		return list.Items[0]
	}
	result, err := rmc.GetProject(ctx, &common.IDOptions{Id: id})
	if err != nil {
		cliLog.Fatal().Err(err).Str("project", id).Msg("Failed to get project")
	}
	return result
}
