//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package platform

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	common "github.com/arangodb-managed/apis/common/v1"
	platform "github.com/arangodb-managed/apis/platform/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/format"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.ListCmd,
		&cobra.Command{
			Use:   "regions",
			Short: "List all regions of the given provider",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				organizationID string
				providerID     string
			}{}
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Optional Identifier of the organization")
			f.StringVarP(&cargs.providerID, "provider-id", "p", cmd.DefaultProvider(), "Identifier of the provider")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				providerID, argsUsed := cmd.OptOption("provider-id", cargs.providerID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				platformc := platform.NewPlatformServiceClient(conn)
				rmc := rm.NewResourceManagerServiceClient(conn)
				ctx := cmd.ContextWithToken()

				var orgID string
				// Fetch organization
				if cargs.organizationID != "" {
					org := selection.MustSelectOrganization(ctx, log, cargs.organizationID, rmc)
					orgID = org.GetId()
				}

				// Fetch provider
				provider := selection.MustSelectProvider(ctx, log, providerID, cargs.organizationID, platformc)

				// Fetch regions in provider
				list, err := platformc.ListRegions(ctx, &platform.ListRegionsRequest{ProviderId: provider.GetId(), OrganizationId: orgID, Options: &common.ListOptions{}})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to list regions")
				}

				// Show result
				fmt.Println(format.RegionList(list.Items, cmd.RootArgs.Format))
			}
		},
	)
}
