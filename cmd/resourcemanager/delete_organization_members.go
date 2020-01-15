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
	// deleteOrgMembersCmd deletes a list of members from an organization
	deleteOrgMembersCmd = &cobra.Command{
		Use:   "group",
		Short: "Add members to group",
		Run:   deleteGroupMembersCmdRun,
	}
	deleteOrgMembersArgs struct {
		organizationID string
		userIDs        []string
	}
)

func init() {
	cmd.DeleteCmd.AddCommand(deleteOrgMembersCmd)

	f := deleteOrgMembersCmd.Flags()
	f.StringVarP(&deleteOrgMembersArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
	deleteOrgMembersArgs.userIDs = *f.StringSliceP("user-ids", "u", []string{}, "A coma separated list of user ids")

}

func deleteGroupMembersCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := deleteOrgMembersArgs
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
		log.Fatal().Str("organization_id", organizationID).Msg("May not delete members from deleted organization.")
	}
	var members *rm.MemberList
	// Create Member list
	for _, u := range cargs.userIDs {
		if ok, err := rmc.IsMemberOfOrganization(ctx, &rm.IsMemberOfOrganizationRequest{UserId: u, OrganizationId: organizationID}); err != nil {
			log.Fatal().Err(err).Msg("Can't determine if user is part of the organization or not.")
		} else if !ok.Member {
			log.Fatal().Str("user-id", u).Str("organization-id", organizationID).Msg("User is not a member of the organization.")
		}
		member := &rm.Member{}
		member.UserId = u
	}
	// Delete members
	_, err = rmc.DeleteOrganizationMembers(ctx, &rm.OrganizationMembersRequest{
		OrganizationId: organizationID,
		Members:        members,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to delete users.")
	}
	fmt.Println("Success!")
}
