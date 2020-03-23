//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
//
// Author Gergely Brautigam
//

package iam

import (
	"fmt"

	"github.com/spf13/cobra"

	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

var (
	// addGroupMembersCmd adds a list of members to a group
	addGroupMembersCmd = &cobra.Command{
		Use:   "members",
		Short: "Add members to group",
		Run:   addGroupMembersCmdRun,
	}
	addGroupMembersArgs struct {
		organizationID string
		groupID        string
		userEmails     *[]string
	}
)

func init() {
	addGroupCmd.AddCommand(addGroupMembersCmd)

	f := addGroupMembersCmd.Flags()
	f.StringVarP(&addGroupMembersArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
	f.StringVarP(&addGroupMembersArgs.groupID, "group-id", "g", cmd.DefaultGroup(), "Identifier of the group to add members to")
	addGroupMembersArgs.userEmails = f.StringSliceP("user-emails", "u", []string{}, "A comma separated list of user email addresses")
}

func addGroupMembersCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := addGroupMembersArgs
	groupID, argsUsed := cmd.OptOption("group-id", cargs.groupID, args, 0)
	organizationID, argsUsed := cmd.OptOption("organiztaion-id", cargs.organizationID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	organization := selection.MustSelectOrganization(ctx, log, organizationID, rmc)
	group := selection.MustSelectGroup(ctx, log, groupID, organization.Id, iamc, rmc)

	log.Info().Msgf("Adding members: %s", cargs.userEmails)
	var userIds []string
	members, err := rmc.ListOrganizationMembers(ctx, &common.ListOptions{ContextId: organization.Id})
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
			log.Fatal().Str("email", e).Msg("User not found or not part of the ogranization")
		} else {
			userIds = append(userIds, id)
		}
	}

	if _, err := iamc.AddGroupMembers(ctx, &iam.GroupMembersRequest{GroupId: group.Id, UserIds: userIds}); err != nil {
		log.Fatal().Err(err).Msg("Failed to add users.")
	}

	fmt.Println("Success!")
}
