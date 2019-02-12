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
	// listGroupsCmd fetches groups of the given organization
	listGroupsCmd = &cobra.Command{
		Use:   "groups",
		Short: "List all groups of the given organization",
		Run:   listGroupsCmdRun,
	}
	listGroupsArgs struct {
		organizationID string
	}
)

func init() {
	listCmd.AddCommand(listGroupsCmd)
	f := listGroupsCmd.Flags()
	f.StringVarP(&listGroupsArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func listGroupsCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	organizationID, argsUsed := optOption("organization-id", listGroupsArgs.organizationID, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch organization
	org := mustSelectOrganization(ctx, organizationID, rmc)

	// Fetch groups in organization
	list, err := iamc.ListGroups(ctx, &common.ListOptions{ContextId: org.GetId()})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to list groups")
	}

	// Show result
	fmt.Println(format.GroupList(list.Items, rootArgs.format))
}
