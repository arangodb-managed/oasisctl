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
	// getRolesCmd fetches roles of the given organization
	getRolesCmd = &cobra.Command{
		Use:   "roles",
		Short: "Get all roles of the given organization",
		Run:   getRolesCmdRun,
	}
	getRolesArgs struct {
		organizationID string
	}
)

func init() {
	getCmd.AddCommand(getRolesCmd)
	f := getRolesCmd.Flags()
	f.StringVarP(&getRolesArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func getRolesCmdRun(cmd *cobra.Command, args []string) {
	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch organization
	org := mustSelectOrganization(ctx, getRolesArgs.organizationID, rmc)

	// Fetch roles in organization
	list, err := iamc.ListRoles(ctx, &common.ListOptions{ContextId: org.GetId()})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to list roles")
	}

	// Show result
	fmt.Println(format.RoleList(list.Items, rootArgs.format))
}
