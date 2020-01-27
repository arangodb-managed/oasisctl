//
// DISCLAIMER
//
// Copyright 2019 ArangoDB GmbH, Cologne, Germany
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

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.ListCmd,
		&cobra.Command{
			Use:   "ipwhitelists",
			Short: "List all IP whitelists of the given project",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				organizationID string
				projectID      string
			}{}
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				projectID, argsUsed := cmd.OptOption("project-id", cargs.projectID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				securityc := security.NewSecurityServiceClient(conn)
				rmc := rm.NewResourceManagerServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Fetch project
				project := selection.MustSelectProject(ctx, log, projectID, cargs.organizationID, rmc)

				// Fetch IP whitelists in project
				list, err := securityc.ListIPWhitelists(ctx, &common.ListOptions{ContextId: project.GetId()})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to list IP whitelists")
				}

				// Show result
				fmt.Println(format.IPWhitelistList(list.Items, cmd.RootArgs.Format))
			}
		},
	)
}
