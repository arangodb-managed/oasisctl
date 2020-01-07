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
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
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
	cmd.ListCmd.AddCommand(listProjectsCmd)
	f := listProjectsCmd.Flags()
	f.StringVarP(&listProjectsArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
}

func listProjectsCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := listProjectsArgs
	organizationID, argsUsed := cmd.OptOption("organization-id", cargs.organizationID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch organization
	org := selection.MustSelectOrganization(ctx, log, organizationID, rmc)

	// Fetch projects in organization
	list, err := rmc.ListProjects(ctx, &common.ListOptions{ContextId: org.GetId()})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list projects")
	}

	// Show result
	fmt.Println(format.ProjectList(list.Items, cmd.RootArgs.Format))
}
