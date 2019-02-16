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

	common "github.com/arangodb-managed/apis/common/v1"
	crypto "github.com/arangodb-managed/apis/crypto/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.DeleteCmd,
		&cobra.Command{
			Use:   "cacertificate",
			Short: "Delete a CA certificate the authenticated user has access to",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				organizationID string
				projectID      string
				cacertID       string
			}{}
			f.StringVarP(&cargs.cacertID, "cacertificate-id", "c", cmd.DefaultCACertificate(), "Identifier of the CA certificate")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")

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

				// Delete CA certificate
				if _, err := cryptoc.DeleteCACertificate(ctx, &common.IDOptions{Id: item.GetId()}); err != nil {
					log.Fatal().Err(err).Msg("Failed to delete CA certificate")
				}

				// Show result
				fmt.Println("Deleted CA certificate!")
			}
		},
	)
}
