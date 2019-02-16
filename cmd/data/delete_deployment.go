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
	"github.com/arangodb-managed/oasis/pkg/selection"
)

var (
	// deleteDeploymentCmd deletes a deployment that the user has access to
	deleteDeploymentCmd = &cobra.Command{
		Use:   "deployment",
		Short: "Delete a deployment the authenticated user has access to",
		Run:   deleteDeploymentCmdRun,
	}
	deleteDeploymentArgs struct {
		organizationID string
		projectID      string
		deploymentID   string
	}
)

func init() {
	cmd.DeleteCmd.AddCommand(deleteDeploymentCmd)
	f := deleteDeploymentCmd.Flags()
	f.StringVarP(&deleteDeploymentArgs.deploymentID, "deployment-id", "d", cmd.DefaultDeployment(), "Identifier of the deployment")
	f.StringVarP(&deleteDeploymentArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
	f.StringVarP(&deleteDeploymentArgs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")
}

func deleteDeploymentCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := deleteDeploymentArgs
	deploymentID, argsUsed := cmd.OptOption("deployment-id", cargs.deploymentID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	datac := data.NewDataServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch deployment
	item := selection.MustSelectDeployment(ctx, log, deploymentID, cargs.projectID, cargs.organizationID, datac, rmc)

	// Delete deployment
	if _, err := datac.DeleteDeployment(ctx, &common.IDOptions{Id: item.GetId()}); err != nil {
		cmd.CLILog.Fatal().Err(err).Msg("Failed to delete deployment")
	}

	// Show result
	fmt.Println("Deleted deployment!")
}
