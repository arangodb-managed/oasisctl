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
	// acceptOrganizationInviteCmd accepts an organization invite that the user has access to
	acceptOrganizationInviteCmd = &cobra.Command{
		Use:   "invite",
		Short: "Accept an organization invite the authenticated user has access to",
		Run:   acceptOrganizationInviteCmdRun,
	}
	acceptOrganizationInviteArgs struct {
		organizationID string
		inviteID       string
	}
)

func init() {
	acceptOrganizationCmd.AddCommand(acceptOrganizationInviteCmd)
	f := acceptOrganizationInviteCmd.Flags()
	f.StringVarP(&acceptOrganizationInviteArgs.inviteID, "invite-id", "i", defaultOrganizationInvite(), "Identifier of the organization invite")
	f.StringVarP(&acceptOrganizationInviteArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func acceptOrganizationInviteCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	inviteID, argsUsed := optOption("invite-id", acceptOrganizationInviteArgs.inviteID, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch invite
	invite := mustSelectOrganizationInvite(ctx, inviteID, acceptOrganizationInviteArgs.organizationID, rmc)

	// Accept invite
	if _, err := rmc.AcceptOrganizationInvite(ctx, invite); err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to accept organization invite")
	}

	// Fetch organization
	orgName := invite.GetOrganizationId()
	if org, err := rmc.GetOrganization(ctx, &common.IDOptions{Id: invite.GetOrganizationId()}); err != nil {
		cliLog.Warn().Err(err).Msg("Failed to get organization")
	} else {
		orgName = org.GetName()
	}

	// Show result
	fmt.Println("Success!")
	fmt.Printf("You are now a member of the '%s' organization.\n", orgName)
}
