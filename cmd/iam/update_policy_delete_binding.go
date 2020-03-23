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
// Author Ewout Prangsma
//

package iam

import (
	"fmt"

	"github.com/spf13/cobra"

	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

var (
	// updatePolicyDeleteBindingCmd deleted a role binding from a policy
	updatePolicyDeleteBindingCmd = &cobra.Command{
		Use:   "binding",
		Short: "Delete a role binding from a policy",
		Run:   updatePolicyDeleteBindingCmdRun,
	}
	updatePolicyDeleteBindingArgs struct {
		url      string
		roleID   string
		userIDs  []string
		groupIDs []string
	}
)

func init() {
	updatePolicyDeleteCmd.AddCommand(updatePolicyDeleteBindingCmd)
	f := updatePolicyDeleteBindingCmd.Flags()
	f.StringVarP(&updatePolicyDeleteBindingArgs.url, "url", "u", cmd.DefaultURL(), "URL of the resource to update the policy for")
	f.StringVarP(&updatePolicyDeleteBindingArgs.roleID, "role-id", "r", cmd.DefaultRole(), "Identifier of the role to delete bind for")
	f.StringSliceVar(&updatePolicyDeleteBindingArgs.userIDs, "user-id", nil, "Identifiers of the users to delete bindings for")
	f.StringSliceVar(&updatePolicyDeleteBindingArgs.groupIDs, "group-id", nil, "Identifiers of the groups to delete bindings for")
}

func updatePolicyDeleteBindingCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := updatePolicyDeleteBindingArgs
	url, argsUsed := cmd.OptOption("url", cargs.url, args, 0)
	roleID, _ := cmd.ReqOption("role-id", cargs.roleID, nil, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)
	if len(cargs.userIDs) == 0 &&
		len(cargs.groupIDs) == 0 {
		log.Fatal().Msg("Provide at least one --user-id or --group-id")
	}

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Parse URL to get organization ID from URL
	resURL, err := rm.ParseResourceURL(cargs.url)
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid resource URL")
	}

	// Get organization ID
	orgID := resURL.OrganizationID()

	// Fetch role
	role := selection.MustSelectRole(ctx, log, roleID, orgID, iamc, rmc)

	// Add role binding
	req := &iam.RoleBindingsRequest{
		ResourceUrl: url,
	}
	for _, uid := range cargs.userIDs {
		// Append users
		item := selection.MustSelectMember(ctx, log, uid, orgID, iamc, rmc)
		req.Bindings = append(req.Bindings, &iam.RoleBinding{
			MemberId: iam.CreateMemberIDFromUserID(item.GetId()),
			RoleId:   role.GetId(),
		})
	}
	for _, gid := range cargs.groupIDs {
		// Append groups
		item := selection.MustSelectGroup(ctx, log, gid, orgID, iamc, rmc)
		req.Bindings = append(req.Bindings, &iam.RoleBinding{
			MemberId: iam.CreateMemberIDFromGroupID(item.GetId()),
			RoleId:   role.GetId(),
		})
	}
	updated, err := iamc.DeleteRoleBindings(ctx, req)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to update policy")
	}

	// Show result
	fmt.Println("Updated policy!")
	fmt.Println(format.Policy(ctx, updated, iamc, cmd.RootArgs.Format))
}
