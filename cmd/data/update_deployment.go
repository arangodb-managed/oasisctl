//
// DISCLAIMER
//
// Copyright 2019 ArangoDB GmbH, Cologne, Germany
//
// Author Ewout Prangsma
//

package data

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	data "github.com/arangodb-managed/apis/data/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.UpdateCmd,
		&cobra.Command{
			Use:   "deployment",
			Short: "Update a deployment the authenticated user has access to",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				deploymentID   string
				organizationID string
				projectID      string
				name           string
				description    string
				ipwhitelistID  string
			}{}
			f.StringVarP(&cargs.deploymentID, "deployment-id", "d", cmd.DefaultDeployment(), "Identifier of the deployment")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")
			f.StringVar(&cargs.name, "name", "", "Name of the deployment")
			f.StringVar(&cargs.description, "description", "", "Description of the deployment")
			f.StringVarP(&cargs.ipwhitelistID, "ipwhitelist-id", "i", cmd.DefaultIPWhitelist(), "Identifier of the IP whitelist to use for the deployment")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				deploymentID, argsUsed := cmd.OptOption("deployment-id", cargs.deploymentID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				datac := data.NewDataServiceClient(conn)
				rmc := rm.NewResourceManagerServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Fetch deployment
				item := selection.MustSelectDeployment(ctx, log, deploymentID, cargs.projectID, cargs.organizationID, datac, rmc)

				// Set changes
				f := c.Flags()
				hasChanges := false
				if f.Changed("name") {
					item.Name = cargs.name
					hasChanges = true
				}
				if f.Changed("description") {
					item.Description = cargs.description
					hasChanges = true
				}
				if f.Changed("ipwhitelist-id") {
					item.IpwhitelistId = cargs.ipwhitelistID
					hasChanges = true
				}
				if !hasChanges {
					fmt.Println("No changes")
				} else {
					// Update deployment
					updated, err := datac.UpdateDeployment(ctx, item)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to update deployment")
					}

					// Show result
					fmt.Println("Updated deployment!")
					fmt.Println(format.Deployment(updated, nil, cmd.RootArgs.Format, false))
				}
			}
		},
	)
}
