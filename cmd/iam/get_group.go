//
// DISCLAIMER
//
// Copyright 2019 ArangoDB GmbH, Cologne, Germany
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
	cmd.GetCmd.AddCommand(getGroupCmd)
	f := getGroupCmd.Flags()
	f.StringVarP(&getGroupArgs.groupID, "group-id", "g", cmd.DefaultGroup(), "Identifier of the group")
	f.StringVarP(&getGroupArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
}

func getGroupCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := getGroupArgs
	groupID, argsUsed := cmd.OptOption("group-id", cargs.groupID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch group
	item := selection.MustSelectGroup(ctx, log, groupID, cargs.organizationID, iamc, rmc)

	// Show result
	fmt.Println(format.Group(item, cmd.RootArgs.Format))
}
