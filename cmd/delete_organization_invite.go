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
)

var (
	// deleteOrganizationInviteCmd deletes an organization invite that the user has access to
	deleteOrganizationInviteCmd = &cobra.Command{
		Use:   "invite",
		Short: "Delete an organization invite the authenticated user has access to",
		Run:   deleteOrganizationInviteCmdRun,
	}
	deleteOrganizationInviteArgs struct {
		organizationID string
		inviteID       string
	}
)

func init() {
	deleteOrganizationCmd.AddCommand(deleteOrganizationInviteCmd)
	f := deleteOrganizationInviteCmd.Flags()
	f.StringVarP(&deleteOrganizationInviteArgs.inviteID, "invite-id", "i", defaultOrganizationInvite(), "Identifier of the organization invite")
	f.StringVarP(&deleteOrganizationInviteArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func deleteOrganizationInviteCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	inviteID, argsUsed := optOption("invite-id", deleteOrganizationInviteArgs.inviteID, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch invite
	item := mustSelectOrganizationInvite(ctx, inviteID, deleteOrganizationInviteArgs.organizationID, rmc)

	// Delete invite
	if _, err := rmc.DeleteOrganizationInvite(ctx, item); err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to delete organization invite")
	}

	// Show result
	fmt.Println("Deleted organization invite!")
}
