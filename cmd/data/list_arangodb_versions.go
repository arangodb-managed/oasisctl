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

	common "github.com/arangodb-managed/apis/common/v1"
	data "github.com/arangodb-managed/apis/data/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.ListArangoDBCmd,
		&cobra.Command{
			Use:   "versions",
			Short: "List all supported ArangoDB versions",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				organizationID string
			}{}
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Optional Identifier of the organization")
			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				cmd.MustCheckNumberOfArgs(args, 0)

				// Connect
				conn := cmd.MustDialAPI()
				datac := data.NewDataServiceClient(conn)
				rmc := rm.NewResourceManagerServiceClient(conn)
				ctx := cmd.ContextWithToken()

				var orgID string
				// Fetch organization
				if cargs.organizationID != "" {
					org := selection.MustSelectOrganization(ctx, log, cargs.organizationID, rmc)
					orgID = org.GetId()
				}

				// Fetch versions
				list, err := datac.ListVersions(ctx, &data.ListVersionsRequest{OrganizationId: orgID, Options: &common.ListOptions{}})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to list versions")
				}

				// Fetch default version
				defaultVersion, err := datac.GetDefaultVersion(ctx, &common.Empty{})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to get default version")
				}

				// Show result
				fmt.Println(format.VersionList(list.Items, defaultVersion, cmd.RootArgs.Format))
			}
		},
	)
}
