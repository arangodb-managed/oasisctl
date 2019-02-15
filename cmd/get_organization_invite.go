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

	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/pkg/format"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

var (
	// getOrganizationInviteCmd fetches an organization invite that the user has access to
	getOrganizationInviteCmd = &cobra.Command{
		Use:   "invite",
		Short: "Get an organization invite the authenticated user has access to",
		Run:   getOrganizationInviteCmdRun,
	}
	getOrganizationInviteArgs struct {
		organizationID string
		inviteID       string
	}
)

func init() {
	getOrganizationCmd.AddCommand(getOrganizationInviteCmd)
	f := getOrganizationInviteCmd.Flags()
	f.StringVarP(&getOrganizationInviteArgs.inviteID, "invite-id", "i", defaultOrganizationInvite(), "Identifier of the organization invite")
	f.StringVarP(&getOrganizationInviteArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func getOrganizationInviteCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	inviteID, argsUsed := optOption("invite-id", getOrganizationInviteArgs.inviteID, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch organization invite
	item := selection.MustSelectOrganizationInvite(ctx, cliLog, inviteID, getOrganizationInviteArgs.organizationID, rmc)

	// Show result
	fmt.Println(format.OrganizationInvite(ctx, item, iamc, rootArgs.format))
}
