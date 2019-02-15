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

	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/pkg/format"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

var (
	// getOrganizationCmd fetches an organization the user is a part of
	getOrganizationCmd = &cobra.Command{
		Use:   "organization",
		Short: "Get an organization the authenticated user is a member of",
		Run:   getOrganizationCmdRun,
	}
	getOrganizationArgs struct {
		organizationID string
	}
)

func init() {
	getCmd.AddCommand(getOrganizationCmd)
	f := getOrganizationCmd.Flags()
	f.StringVarP(&getOrganizationArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func getOrganizationCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	organizationID, argsUsed := optOption("organization-id", getOrganizationArgs.organizationID, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch organization
	item := selection.MustSelectOrganization(ctx, cliLog, organizationID, rmc)

	// Show result
	fmt.Println(format.Organization(item, rootArgs.format))
}
