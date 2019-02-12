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
	// getOrganizationInvitesCmd fetches the invites of an organization the user is a part of
	getOrganizationInvitesCmd = &cobra.Command{
		Use:   "invites",
		Short: "Get invites of an organization the authenticated user is a member of",
		Run:   getOrganizationInvitesCmdRun,
	}
	getOrganizationInvitesArgs struct {
		organizationID string
	}
)

func init() {
	getOrganizationCmd.AddCommand(getOrganizationInvitesCmd)
	f := getOrganizationInvitesCmd.Flags()
	f.StringVarP(&getOrganizationInvitesArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func getOrganizationInvitesCmdRun(cmd *cobra.Command, args []string) {
	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch organization
	org := mustSelectOrganization(ctx, getOrganizationInvitesArgs.organizationID, rmc)

	list, err := rmc.ListOrganizationInvites(ctx, &common.ListOptions{ContextId: org.GetId()})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to list organization invites")
	}

	// Show result
	fmt.Println(format.OrganizationInviteList(ctx, list.GetItems(), iamc, rootArgs.format))
}
