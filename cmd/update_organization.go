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
)

var (
	// updateOrganizationCmd updates an organization that the user has access to
	updateOrganizationCmd = &cobra.Command{
		Use:   "organization",
		Short: "Update an organization the authenticated user has access to",
		Run:   updateOrganizationCmdRun,
	}
	updateOrganizationArgs struct {
		organizationID string
		name           string
		description    string
	}
)

func init() {
	updateCmd.AddCommand(updateOrganizationCmd)
	f := updateOrganizationCmd.Flags()
	f.StringVarP(&updateOrganizationArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
	f.StringVarP(&updateOrganizationArgs.name, "name", "n", "", "Name of the organization")
	f.StringVarP(&updateOrganizationArgs.description, "description", "d", "", "Description of the organization")
}

func updateOrganizationCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	organizationID, argsUsed := optOption("organization-id", updateOrganizationArgs.organizationID, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch organization
	item := mustSelectOrganization(ctx, organizationID, rmc)

	// Set changes
	f := cmd.Flags()
	hasChanges := false
	if f.Changed("name") {
		item.Name = updateOrganizationArgs.name
		hasChanges = true
	}
	if f.Changed("description") {
		item.Description = updateOrganizationArgs.description
		hasChanges = true
	}
	if !hasChanges {
		fmt.Println("No changes")
	} else {
		// Update project
		updated, err := rmc.UpdateOrganization(ctx, item)
		if err != nil {
			cliLog.Fatal().Err(err).Msg("Failed to update organization")
		}

		// Show result
		fmt.Println("Updated organization!")
		fmt.Println(format.Organization(updated, rootArgs.format))
	}
}
