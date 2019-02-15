//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package rm

import (
	"github.com/arangodb-managed/oasis/cmd"
	"fmt"

	"github.com/spf13/cobra"

	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
	"github.com/arangodb-managed/oasis/pkg/selection"
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
	f.StringVarP(&deleteOrganizationInviteArgs.inviteID, "invite-id", "i", cmd.DefaultOrganizationInvite(), "Identifier of the organization invite")
	f.StringVarP(&deleteOrganizationInviteArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
}

func deleteOrganizationInviteCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	inviteID, argsUsed := cmd.OptOption("invite-id", deleteOrganizationInviteArgs.inviteID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch invite
	item := selection.MustSelectOrganizationInvite(ctx, cmd.CLILog, inviteID, deleteOrganizationInviteArgs.organizationID, rmc)

	// Delete invite
	if _, err := rmc.DeleteOrganizationInvite(ctx, item); err != nil {
		cmd.CLILog.Fatal().Err(err).Msg("Failed to delete organization invite")
	}

	// Show result
	fmt.Println("Deleted organization invite!")
}
