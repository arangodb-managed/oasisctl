//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package security

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
	security "github.com/arangodb-managed/apis/security/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.DeleteCmd,
		&cobra.Command{
			Use:   "ipwhitelist",
			Short: "Delete an IP whitelist the authenticated user has access to",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				organizationID string
				projectID      string
				ipwhitelistID  string
			}{}
			f.StringVarP(&cargs.ipwhitelistID, "ipwhitelist-id", "i", cmd.DefaultIPWhitelist(), "Identifier of the IP whitelist")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				ipwhitelistID, argsUsed := cmd.OptOption("ipwhitelist-id", cargs.ipwhitelistID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				securityc := security.NewSecurityServiceClient(conn)
				rmc := rm.NewResourceManagerServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Fetch IP whitelist
				item := selection.MustSelectIPWhitelist(ctx, log, ipwhitelistID, cargs.projectID, cargs.organizationID, securityc, rmc)

				// Delete IP whitelist
				if _, err := securityc.DeleteIPWhitelist(ctx, &common.IDOptions{Id: item.GetId()}); err != nil {
					log.Fatal().Err(err).Msg("Failed to delete IP whitelist")
				}

				// Show result
				fmt.Println("Deleted IP whitelist!")
			}
		},
	)
}
