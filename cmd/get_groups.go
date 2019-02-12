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
	// getGroupsCmd fetches groups of the given organization
	getGroupsCmd = &cobra.Command{
		Use:   "groups",
		Short: "Get all groups of the given organization",
		Run:   getGroupsCmdRun,
	}
	getGroupsArgs struct {
		organizationID string
	}
)

func init() {
	getCmd.AddCommand(getGroupsCmd)
	f := getGroupsCmd.Flags()
	f.StringVarP(&getGroupsArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func getGroupsCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	mustCheckNumberOfArgs(args, 0)

	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch organization
	org := mustSelectOrganization(ctx, getGroupsArgs.organizationID, rmc)

	// Fetch groups in organization
	list, err := iamc.ListGroups(ctx, &common.ListOptions{ContextId: org.GetId()})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to list groups")
	}

	// Show result
	fmt.Println(format.GroupList(list.Items, rootArgs.format))
}
