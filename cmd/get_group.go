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
	// getGroupCmd fetches a group that the user has access to
	getGroupCmd = &cobra.Command{
		Use:   "group",
		Short: "Get a group the authenticated user has access to",
		Run:   getGroupCmdRun,
	}
	getGroupArgs struct {
		groupID        string
		organizationID string
	}
)

func init() {
	getCmd.AddCommand(getGroupCmd)
	f := getGroupCmd.Flags()
	f.StringVarP(&getGroupArgs.groupID, "group-id", "g", defaultGroup(), "Identifier of the group")
	f.StringVarP(&getGroupArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func getGroupCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	groupID, argsUsed := optOption("group-id", getGroupArgs.groupID, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch group
	item := selection.MustSelectGroup(ctx, cliLog, groupID, getGroupArgs.organizationID, iamc, rmc)

	// Show result
	fmt.Println(format.Group(item, rootArgs.format))
}
