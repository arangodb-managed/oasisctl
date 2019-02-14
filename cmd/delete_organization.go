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
	deleteCmd.AddCommand(deleteOrganizationCmd)
	f := deleteOrganizationCmd.Flags()
	f.StringVarP(&deleteOrganizationArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func deleteOrganizationCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	organizationID, argsUsed := optOption("organization-id", deleteOrganizationArgs.organizationID, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch organization
	item := mustSelectOrganization(ctx, organizationID, rmc)

	// Delete project
	if _, err := rmc.DeleteOrganization(ctx, &common.IDOptions{Id: item.GetId()}); err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to delete organization")
	}

	// Show result
	fmt.Println("Deleted organization!")
}
