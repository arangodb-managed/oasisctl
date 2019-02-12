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
	// createOrganizationCmd creates a new organization
	createOrganizationCmd = &cobra.Command{
		Use:   "organization",
		Short: "Create a new organization",
		Run:   createOrganizationCmdRun,
	}
	createOrganizationArgs struct {
		name        string
		description string
	}
)

func init() {
	createCmd.AddCommand(createOrganizationCmd)

	f := createOrganizationCmd.Flags()
	f.StringVarP(&createOrganizationArgs.name, "name", "n", "", "Name of the organization")
	f.StringVarP(&createOrganizationArgs.description, "description", "d", "", "Description of the organization")
}

func createOrganizationCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	name, argsUsed := reqOption("name", createOrganizationArgs.name, args, 0)
	description := createOrganizationArgs.description
	mustCheckNumberOfArgs(args, argsUsed)

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
