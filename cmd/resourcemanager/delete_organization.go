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
	"github.com/arangodb-managed/oasis/cmd"

	"github.com/spf13/cobra"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

var (
	// deleteOrganizationCmd deletes an organization that the user has access to
	deleteOrganizationCmd = &cobra.Command{
		Use:   "organization",
		Short: "Delete an organization the authenticated user has access to",
		Run:   deleteOrganizationCmdRun,
	}
	deleteOrganizationArgs struct {
		organizationID string
	}
)

func init() {
	cmd.DeleteCmd.AddCommand(deleteOrganizationCmd)
	f := deleteOrganizationCmd.Flags()
	f.StringVarP(&deleteOrganizationArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
}

func deleteOrganizationCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	organizationID, argsUsed := cmd.OptOption("organization-id", deleteOrganizationArgs.organizationID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch organization
	item := selection.MustSelectOrganization(ctx, cmd.CLILog, organizationID, rmc)

	// Delete project
	if _, err := rmc.DeleteOrganization(ctx, &common.IDOptions{Id: item.GetId()}); err != nil {
		cmd.CLILog.Fatal().Err(err).Msg("Failed to delete organization")
	}

	// Show result
	fmt.Println("Deleted organization!")
}
