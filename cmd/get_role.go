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
	item := selection.MustSelectRole(ctx, cliLog, roleID, getRoleArgs.organizationID, iamc, rmc)

	// Show result
	fmt.Println(format.Role(item, rootArgs.format))
}
