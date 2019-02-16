//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package rm

import (
	"fmt"
	"github.com/arangodb-managed/oasis/cmd"

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
	cmd.CreateCmd.AddCommand(createOrganizationCmd)

	f := createOrganizationCmd.Flags()
	f.StringVar(&createOrganizationArgs.name, "name", "", "Name of the organization")
	f.StringVar(&createOrganizationArgs.description, "description", "", "Description of the organization")
}

func createOrganizationCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	name, argsUsed := cmd.ReqOption("name", createOrganizationArgs.name, args, 0)
	description := createOrganizationArgs.description
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Create organization
	result, err := rmc.CreateOrganization(ctx, &rm.Organization{
		Name:        name,
		Description: description,
	})
	if err != nil {
		cmd.CLILog.Fatal().Err(err).Msg("Failed to create organization")
	}

	// Show result
	fmt.Println("Success!")
	fmt.Println(format.Organization(result, cmd.RootArgs.Format))
}
