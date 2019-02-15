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
	"github.com/arangodb-managed/oasis/pkg/selection"
)

var (
	// listProjectsCmd fetches projects of the given organization
	listProjectsCmd = &cobra.Command{
		Use:   "projects",
		Short: "List all projects of the given organization",
		Run:   listProjectsCmdRun,
	}
	listProjectsArgs struct {
		organizationID string
	}
)

func init() {
	listCmd.AddCommand(listProjectsCmd)
	f := listProjectsCmd.Flags()
	f.StringVarP(&listProjectsArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func listProjectsCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	organizationID, argsUsed := optOption("organization-id", listProjectsArgs.organizationID, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch organization
	org := selection.MustSelectOrganization(ctx, cliLog, organizationID, rmc)

	// Fetch projects in organization
	list, err := rmc.ListProjects(ctx, &common.ListOptions{ContextId: org.GetId()})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to list projects")
	}

	// Show result
	fmt.Println(format.ProjectList(list.Items, rootArgs.format))
}
