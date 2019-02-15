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
	// createOrganizationInvite creates a new organization invite
	createOrganizationInviteCmd = &cobra.Command{
		Use:   "invite",
		Short: "Create a new invite to an organization",
		Run:   createOrganizationInviteCmdRun,
	}
	createOrganizationInviteArgs struct {
		email          string
		organizationID string
	}
)

func init() {
	createOrganizationCmd.AddCommand(createOrganizationInviteCmd)

	f := createOrganizationInviteCmd.Flags()
	f.StringVarP(&createOrganizationInviteArgs.email, "email", "e", "", "Email address of the person to invite")
	f.StringVarP(&createOrganizationInviteArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization to create the invite in")
}

func createOrganizationInviteCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	email, argsUsed := reqOption("email", createOrganizationInviteArgs.email, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch organization
	org := selection.MustSelectOrganization(ctx, cliLog, createOrganizationInviteArgs.organizationID, rmc)

	// Create invite
	result, err := rmc.CreateOrganizationInvite(ctx, &rm.OrganizationInvite{
		OrganizationId: org.GetId(),
		Email:          email,
	})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to create organization invite")
	}

	// Show result
	fmt.Println("Success!")
	fmt.Println(format.OrganizationInvite(ctx, result, iamc, rootArgs.format))
}
