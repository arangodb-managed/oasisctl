//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/pkg/format"
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
	item := mustSelectOrganizationInvite(ctx, inviteID, getOrganizationInviteArgs.organizationID, rmc)

	// Show result
	fmt.Println(format.OrganizationInvite(ctx, item, iamc, rootArgs.format))
}

// mustSelectOrganizationInvite fetches the organization invite with given ID.
// If no ID is specified, all invites are fetched from the selected organization
// and if the list is exactly 1 long, that invite is returned.
func mustSelectOrganizationInvite(ctx context.Context, id, orgID string, rmc rm.ResourceManagerServiceClient) *rm.OrganizationInvite {
	if id == "" {
		org := mustSelectOrganization(ctx, orgID, rmc)
		list, err := rmc.ListOrganizationInvites(ctx, &common.ListOptions{ContextId: org.GetId()})
		if err != nil {
			cliLog.Fatal().Err(err).Msg("Failed to list organization invites")
		}
		if len(list.Items) != 1 {
			cliLog.Fatal().Err(err).Msgf("You have access to %d organization invites. Please specify one explicitly.", len(list.Items))
		}
		return list.Items[0]
	}
	result, err := rmc.GetOrganizationInvite(ctx, &common.IDOptions{Id: id})
	if err != nil {
		cliLog.Fatal().Err(err).Str("organization-invite", id).Msg("Failed to get organization invite")
	}
	return result
}
