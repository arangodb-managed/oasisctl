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

	"github.com/spf13/cobra"

	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
)

var (
	// deleteOrgMembersCmd deletes a list of members from an organization
	deleteOrgMembersCmd = &cobra.Command{
		Use:   "members",
		Short: "Delete members from organization",
		Run:   deleteOrgMembersCmdRun,
	}
	deleteOrgMembersArgs struct {
		organizationID string
		userEmails     *[]string
	}
)

func init() {
	deleteOrganizationCmd.AddCommand(deleteOrgMembersCmd)

	f := deleteOrgMembersCmd.Flags()
	f.StringVarP(&deleteOrgMembersArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
	deleteOrgMembersArgs.userEmails = f.StringSliceP("user-emails", "u", []string{}, "A comma separated list of user email addresses")

}

func deleteOrgMembersCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := deleteOrgMembersArgs
	organizationID, argsUsed := cmd.OptOption("organization-id", cargs.organizationID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	log.Info().Msgf("Deleting members: %s", cargs.userEmails)
	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	org, err := rmc.GetOrganization(ctx, &common.IDOptions{Id: organizationID})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get organization.")
	}
	if org.IsDeleted {
		log.Fatal().Str("organization_id", organizationID).Msg("May not delete members from deleted organization.")
	}

	membersToDelete := &rm.MemberList{Items: make([]*rm.Member, 0)}
	members, err := rmc.ListOrganizationMembers(ctx, &common.ListOptions{ContextId: organizationID})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list organization members.")
	}
	emailIDMap := make(map[string]string)
	for _, u := range members.Items {
		user, err := iamc.GetUser(ctx, &common.IDOptions{Id: u.UserId})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get user")
		}
		emailIDMap[user.Email] = user.Id
	}

	for _, e := range *cargs.userEmails {
		if id, ok := emailIDMap[e]; !ok {
			log.Fatal().Str("email", e).Str("organization-id", organizationID).Msg("User is not a member of the organization.")
		} else {
			membersToDelete.Items = append(membersToDelete.Items, &rm.Member{UserId: id})
		}
	}

	if _, err = rmc.DeleteOrganizationMembers(ctx, &rm.OrganizationMembersRequest{
		OrganizationId: organizationID,
		Members:        membersToDelete,
	}); err != nil {
		log.Fatal().Err(err).Msg("Failed to delete users.")
	}

	fmt.Println("Success!")
}
