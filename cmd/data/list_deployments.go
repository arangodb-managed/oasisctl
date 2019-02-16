//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package data

import (
	"fmt"

	"github.com/spf13/cobra"

	common "github.com/arangodb-managed/apis/common/v1"
	data "github.com/arangodb-managed/apis/data/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/format"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

var (
	// listDeploymentsCmd fetches deployments of the given project
	listDeploymentsCmd = &cobra.Command{
		Use:   "deployments",
		Short: "List all deployments of the given project",
		Run:   listDeploymentsCmdRun,
	}
	listDeploymentsArgs struct {
		organizationID string
		projectID      string
	}
)

func init() {
	cmd.ListCmd.AddCommand(listDeploymentsCmd)
	f := listDeploymentsCmd.Flags()
	f.StringVarP(&listDeploymentsArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
	f.StringVarP(&listDeploymentsArgs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")
}

func listDeploymentsCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := listDeploymentsArgs
	projectID, argsUsed := cmd.OptOption("project-id", cargs.projectID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	datac := data.NewDataServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch project
	project := selection.MustSelectProject(ctx, log, projectID, cargs.organizationID, rmc)

	// Fetch deployments in project
	list, err := datac.ListDeployments(ctx, &common.ListOptions{ContextId: project.GetId()})
	if err != nil {
		cmd.CLILog.Fatal().Err(err).Msg("Failed to list deployments")
	}

	// Show result
	fmt.Println(format.DeploymentList(list.Items, cmd.RootArgs.Format))
}