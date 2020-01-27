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

	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
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
	cmd.ListCmd.AddCommand(listGroupsCmd)
	f := listGroupsCmd.Flags()
	f.StringVarP(&listGroupsArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
}

func listGroupsCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := listGroupsArgs
	organizationID, argsUsed := cmd.OptOption("organization-id", cargs.organizationID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch organization
	org := selection.MustSelectOrganization(ctx, log, organizationID, rmc)

	// Fetch groups in organization
	list, err := iamc.ListGroups(ctx, &common.ListOptions{ContextId: org.GetId()})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list groups")
	}

	// Show result
	fmt.Println(format.GroupList(list.Items, cmd.RootArgs.Format))
}
