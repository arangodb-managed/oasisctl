//
// DISCLAIMER
//
// Copyright 2019 ArangoDB GmbH, Cologne, Germany
//
// Author Ewout Prangsma
//

package rm

import (
	"fmt"

	"github.com/spf13/cobra"

	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
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
	f.StringVarP(&listOrganizationInvitesArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
}

func listOrganizationInvitesCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := listOrganizationInvitesArgs
	organizationID, argsUsed := cmd.OptOption("organization-id", cargs.organizationID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch organization
	org := selection.MustSelectOrganization(ctx, log, organizationID, rmc)

	list, err := rmc.ListOrganizationInvites(ctx, &common.ListOptions{ContextId: org.GetId()})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list organization invites")
	}

	// Show result
	fmt.Println(format.OrganizationInviteList(ctx, list.GetItems(), iamc, cmd.RootArgs.Format))
}
