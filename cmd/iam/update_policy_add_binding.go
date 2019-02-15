//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package iam

import (
	"fmt"

	"github.com/spf13/cobra"

	iam "github.com/arangodb-managed/apis/iam/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/format"
)

var (
	// updatePolicyAddBindingCmd adds a role binding to a policy
	updatePolicyAddBindingCmd = &cobra.Command{
		Use:   "binding",
		Short: "Add a role binding to a policy",
		Run:   updatePolicyAddBindingCmdRun,
	}
	updatePolicyAddBindingArgs struct {
		url      string
		roleID   string
		userIDs  []string
		groupIDs []string
	}
)

func init() {
	updatePolicyAddCmd.AddCommand(updatePolicyAddBindingCmd)
	f := updatePolicyAddBindingCmd.Flags()
	f.StringVarP(&updatePolicyAddBindingArgs.url, "url", "u", cmd.DefaultURL(), "URL of the resource to update the policy for")
	f.StringVarP(&updatePolicyAddBindingArgs.roleID, "role-id", "r", cmd.DefaultRole(), "Identifier of the role to bind to")
	f.StringSliceVar(&updatePolicyAddBindingArgs.userIDs, "user-id", nil, "Identifiers of the users to add bindings for")
	f.StringSliceVar(&updatePolicyAddBindingArgs.groupIDs, "group-id", nil, "Identifiers of the groups to add bindings for")
}

func updatePolicyAddBindingCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	url, argsUsed := cmd.OptOption("url", updatePolicyAddBindingArgs.url, args, 0)
	roleID, _ := cmd.ReqOption("role-id", updatePolicyAddBindingArgs.roleID, nil, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)
	if len(updatePolicyAddBindingArgs.userIDs) == 0 &&
		len(updatePolicyAddBindingArgs.groupIDs) == 0 {
		log.Fatal().Msg("Provide at least one --user-id or --group-id")
	}

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Add role binding
	req := &iam.RoleBindingsRequest{
		ResourceUrl: url,
	}
	for _, uid := range updatePolicyAddBindingArgs.userIDs {
		req.Bindings = append(req.Bindings, &iam.RoleBinding{
			MemberId: iam.CreateMemberIDFromUserID(uid),
			RoleId:   roleID,
		})
	}
	for _, gid := range updatePolicyAddBindingArgs.groupIDs {
		req.Bindings = append(req.Bindings, &iam.RoleBinding{
			MemberId: iam.CreateMemberIDFromGroupID(gid),
			RoleId:   roleID,
		})
	}
	updated, err := iamc.AddRoleBindings(ctx, req)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to update policy")
	}

	// Show result
	fmt.Println("Updated policy!")
	fmt.Println(format.Policy(ctx, updated, iamc, cmd.RootArgs.Format))
}
