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
	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/pkg/format"
)

var (
	// listOrganizationInvitesCmd fetches the invites of an organization the user is a part of
	listOrganizationInvitesCmd = &cobra.Command{
		Use:   "invites",
		Short: "List invites of an organization the authenticated user is a member of",
		Run:   listOrganizationInvitesCmdRun,
	}
	listOrganizationInvitesArgs struct {
		organizationID string
	}
)

func init() {
	listOrganizationCmd.AddCommand(listOrganizationInvitesCmd)
	f := listOrganizationInvitesCmd.Flags()
	f.StringVarP(&listOrganizationInvitesArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func listOrganizationInvitesCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	organizationID, argsUsed := optOption("organization-id", listOrganizationInvitesArgs.organizationID, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch organization
	org := mustSelectOrganization(ctx, organizationID, rmc)

	list, err := rmc.ListOrganizationInvites(ctx, &common.ListOptions{ContextId: org.GetId()})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to list organization invites")
	}

	// Show result
	fmt.Println(format.OrganizationInviteList(ctx, list.GetItems(), iamc, rootArgs.format))
}
