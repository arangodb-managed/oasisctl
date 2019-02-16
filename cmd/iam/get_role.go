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
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/cmd"
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
	cmd.GetCmd.AddCommand(getRoleCmd)
	f := getRoleCmd.Flags()
	f.StringVarP(&getRoleArgs.roleID, "role-id", "r", cmd.DefaultRole(), "Identifier of the role")
	f.StringVarP(&getRoleArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
}

func getRoleCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := getRoleArgs
	roleID, argsUsed := cmd.OptOption("role-id", cargs.roleID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch role
	item := selection.MustSelectRole(ctx, log, roleID, cargs.organizationID, iamc, rmc)

	// Show result
	fmt.Println(format.Role(item, cmd.RootArgs.Format))
}
