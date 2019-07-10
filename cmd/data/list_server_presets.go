//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Robert Stam
//

package data

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	data "github.com/arangodb-managed/apis/data/v1"
	platform "github.com/arangodb-managed/apis/platform/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/format"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.ListServersCmd,
		&cobra.Command{
			Use:   "presets",
			Short: "List presets for servers in a project for a specific region",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				regionID       string
				projectID      string
				organizationID string
				providerID     string
			}{}
			f.StringVarP(&cargs.regionID, "region-id", "r", cmd.DefaultRegion(), "Identifier of the region")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
			f.StringVar(&cargs.providerID, "provider-id", cmd.DefaultProvider(), "Identifier of the provider")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				regionID, argsUsed := cmd.OptOption("region-id", cargs.regionID, args, 0)
				projectID, argsUsed := cmd.OptOption("project-id", cargs.projectID, args, 1)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				datac := data.NewDataServiceClient(conn)
				platformc := platform.NewPlatformServiceClient(conn)
				rmc := rm.NewResourceManagerServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Select project
				project := selection.MustSelectProject(ctx, log, projectID, cargs.organizationID, rmc)

				// Selection region
				region := selection.MustSelectRegion(ctx, log, regionID, cargs.providerID, project.OrganizationId, platformc)

				// Fetch presets
				list, err := datac.ListServersSpecPresets(ctx, &data.ServersSpecPresetsRequest{
					ProjectId: project.GetId(),
					RegionId:  region.GetId(),
				})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to list servers presets")
				}

				// Show result
				fmt.Println(format.ServersSpecPresetList(list.Items, cmd.RootArgs.Format))
			}
		},
	)
}
