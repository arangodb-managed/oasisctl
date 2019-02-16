//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package crypto

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	crypto "github.com/arangodb-managed/apis/crypto/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/format"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.UpdateCmd,
		&cobra.Command{
			Use:   "cacertificate",
			Short: "Update a CA certificate the authenticated user has access to",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				cacertID       string
				organizationID string
				projectID      string
				name           string
				description    string
			}{}
			f.StringVarP(&cargs.cacertID, "cacertificate-id", "c", cmd.DefaultCACertificate(), "Identifier of the CA certificate")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")
			f.StringVarP(&cargs.name, "name", "n", "", "Name of the CA certificate")
			f.StringVarP(&cargs.description, "description", "d", "", "Description of the CA certificate")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				cacertID, argsUsed := cmd.OptOption("cacertificate-id", cargs.cacertID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				cryptoc := crypto.NewCryptoServiceClient(conn)
				rmc := rm.NewResourceManagerServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Fetch CA certificate
				item := selection.MustSelectCACertificate(ctx, log, cacertID, cargs.projectID, cargs.organizationID, cryptoc, rmc)

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
				if !hasChanges {
					fmt.Println("No changes")
				} else {
					// Update CA certificate
					updated, err := cryptoc.UpdateCACertificate(ctx, item)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to update CA certificate")
					}

					// Show result
					fmt.Println("Updated CA certificate!")
					fmt.Println(format.CACertificate(updated, cmd.RootArgs.Format))
				}
			}
		},
	)
}
