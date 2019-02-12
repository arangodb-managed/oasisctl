//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/pkg/format"
)

var (
	// getRoleCmd fetches a group that the user has access to
	getRoleCmd = &cobra.Command{
		Use:   "role",
		Short: "Get a role the authenticated user has access to",
		Run:   getRoleCmdRun,
	}
	getRoleArgs struct {
		roleID         string
		organizationID string
	}
)

func init() {
	getCmd.AddCommand(getRoleCmd)
	f := getRoleCmd.Flags()
	f.StringVarP(&getRoleArgs.roleID, "role-id", "r", defaultRole(), "Identifier of the role")
	f.StringVarP(&getRoleArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func getRoleCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	roleID, argsUsed := optOption("role-id", getRoleArgs.roleID, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch role
	item := mustSelectRole(ctx, roleID, getRoleArgs.organizationID, iamc, rmc)

	// Show result
	fmt.Println(format.Role(item, rootArgs.format))
}

// mustSelectRole fetches the role with given ID.
// If no ID is specified, all roles are fetched from the selected organization
// and if the list is exactly 1 long, that role is returned.
func mustSelectRole(ctx context.Context, id, orgID string, iamc iam.IAMServiceClient, rmc rm.ResourceManagerServiceClient) *iam.Role {
	if id == "" {
		org := mustSelectOrganization(ctx, orgID, rmc)
		list, err := iamc.ListRoles(ctx, &common.ListOptions{ContextId: org.GetId()})
		if err != nil {
			cliLog.Fatal().Err(err).Msg("Failed to list roles")
		}
		if len(list.Items) != 1 {
			cliLog.Fatal().Err(err).Msg("You have access to %d roles. Please specify one explicitly.")
		}
		return list.Items[0]
	}
	result, err := iamc.GetRole(ctx, &common.IDOptions{Id: id})
	if err != nil {
		cliLog.Fatal().Err(err).Str("group", id).Msg("Failed to get role")
	}
	return result
}
