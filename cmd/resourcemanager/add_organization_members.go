//
// DISCLAIMER
//
// Copyright 2020 ArangoDB Inc, Cologne, Germany
//
// Author Gergely Brautigam
//

package rm

import (
	"fmt"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/spf13/cobra"

	"github.com/arangodb-managed/oasisctl/cmd"
)

var (
	// addOrgMembersCmd adds a list of members to an organization
	addOrgMembersCmd = &cobra.Command{
		Use:   "organization",
		Short: "Add members to group",
		Run:   addOrgMembersCmdRun,
	}
	addOrgMembersArgs struct {
		organizationID string
		userIDs        []string
	}
)

func init() {
	cmd.CreateCmd.AddCommand(addOrgMembersCmd)

	f := addOrgMembersCmd.Flags()
	f.StringVarP(&addOrgMembersArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
	addOrgMembersArgs.userIDs = *f.StringSliceP("user-ids", "u", []string{}, "A coma separated list of user ids")

}

func addOrgMembersCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := addOrgMembersArgs
	organizationID, argsUsed := cmd.OptOption("organization-id", cargs.organizationID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	org, err := rmc.GetOrganization(ctx, &common.IDOptions{Id: organizationID})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get organization.")
	}
	if org.IsDeleted {
		log.Fatal().Str("organization_id", organizationID).Msg("May not add member to deleted organization.")
	}
	var members *rm.MemberList
	// Create Member list
	for _, u := range cargs.userIDs {
		member := &rm.Member{}
		member.UserId = u
	}
	// Add members
	_, err = rmc.AddOrganizationMembers(ctx, &rm.OrganizationMembersRequest{
		OrganizationId: organizationID,
		Members:        members,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to add users.")
	}
	fmt.Println("Success!")
}
