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

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/pkg/format"
)

var (
	// getProjectsCmd fetches projects of the given organization
	getProjectsCmd = &cobra.Command{
		Use:   "projects",
		Short: "Get all projects of the given organization",
		Run:   getProjectsCmdRun,
	}
	getProjectsArgs struct {
		organizationID string
	}
)

func init() {
	getCmd.AddCommand(getProjectsCmd)
	f := getProjectsCmd.Flags()
	f.StringVarP(&getProjectsArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func getProjectsCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	mustCheckNumberOfArgs(args, 0)

	// Connect
	conn := mustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch organization
	org := mustSelectOrganization(ctx, getProjectsArgs.organizationID, rmc)

	// Fetch projects in organization
	list, err := rmc.ListProjects(ctx, &common.ListOptions{ContextId: org.GetId()})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to list projects")
	}

	// Show result
	fmt.Println(format.ProjectList(list.Items, rootArgs.format))
}
