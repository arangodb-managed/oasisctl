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

	"github.com/arangodb-managed/adbcloud/pkg/format"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

var (
	// createOrganizationCmd creates a new organization
	createOrganizationCmd = &cobra.Command{
		Use:   "organization",
		Short: "Create a new organization",
		Run:   createOrganizationCmdRun,
	}
	createOrganizationOptions struct {
		name        string
		description string
	}
)

func init() {
	createCmd.AddCommand(createOrganizationCmd)

	f := createOrganizationCmd.Flags()
	f.StringVar(&createOrganizationOptions.name, "name", "", "Name of the organization")
	f.StringVar(&createOrganizationOptions.description, "description", "", "Description of the organization")
}

func createOrganizationCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	name := reqOption("name", createOrganizationOptions.name)
	description := createOrganizationOptions.description

	// Connect
	conn := mustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Create organization
	result, err := rmc.CreateOrganization(ctx, &rm.Organization{
		Name:        name,
		Description: description,
	})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to create organization")
	}

	// Show result
	fmt.Println("Success!")
	fmt.Println(format.Organization(result, rootArgs.format))
}
