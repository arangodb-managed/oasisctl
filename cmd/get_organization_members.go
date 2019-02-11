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
	// getOrganizationMembersCmd fetches the members of an organization the user is a part of
	getOrganizationMembersCmd = &cobra.Command{
		Use:   "members",
		Short: "Get members of an organization the authenticated user is a member of",
		Run:   getOrganizationMembersCmdRun,
	}
	getOrganizationMembersArgs struct {
		organizationID string
	}
)

func init() {
	getOrganizationCmd.AddCommand(getOrganizationMembersCmd)
	f := getOrganizationMembersCmd.Flags()
	f.StringVarP(&getOrganizationMembersArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func getOrganizationMembersCmdRun(cmd *cobra.Command, args []string) {
	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch organization
	org := mustSelectOrganization(ctx, getOrganizationMembersArgs.organizationID, rmc)

	list, err := rmc.ListOrganizationMembers(ctx, &common.ListOptions{ContextId: org.GetId()})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to list organization members")
	}

	// Show result
	fmt.Println(format.OrganizationMemberList(ctx, list.GetItems(), iamc, rootArgs.format))
}
