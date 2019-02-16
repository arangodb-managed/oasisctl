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

	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

var (
	// rejectOrganizationInviteCmd rejects an organization invite that the user has access to
	rejectOrganizationInviteCmd = &cobra.Command{
		Use:   "invite",
		Short: "Reject an organization invite the authenticated user has access to",
		Run:   rejectOrganizationInviteCmdRun,
	}
	rejectOrganizationInviteArgs struct {
		organizationID string
		inviteID       string
	}
)

func init() {
	rejectOrganizationCmd.AddCommand(rejectOrganizationInviteCmd)
	f := rejectOrganizationInviteCmd.Flags()
	f.StringVarP(&rejectOrganizationInviteArgs.inviteID, "invite-id", "i", cmd.DefaultOrganizationInvite(), "Identifier of the organization invite")
	f.StringVarP(&rejectOrganizationInviteArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
}

func rejectOrganizationInviteCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := rejectOrganizationInviteArgs
	inviteID, argsUsed := cmd.OptOption("invite-id", cargs.inviteID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch invite
	invite := selection.MustSelectOrganizationInvite(ctx, log, inviteID, cargs.organizationID, rmc)

	// Reject invite
	if _, err := rmc.RejectOrganizationInvite(ctx, invite); err != nil {
		log.Fatal().Err(err).Msg("Failed to reject organization invite")
	}

	// Show result
	fmt.Println("Success!")
	fmt.Println("You have rejected the invite.")
}
