//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
//
// Author Brautigam Gergely
//

package data

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	data "github.com/arangodb-managed/apis/data/v1"
	platform "github.com/arangodb-managed/apis/platform/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.ListCmd,
		&cobra.Command{
			Use:   "nodesizes",
			Short: "List node sizes",
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

				// Fetch node sizes
				list, err := datac.ListNodeSizes(ctx, &data.NodeSizesRequest{
					ProjectId: project.GetId(),
					RegionId:  region.GetId(),
				})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to list node sizes")
				}

				// Show result
				fmt.Println(format.NodeSizeList(list.Items, cmd.RootArgs.Format))
			}
		},
	)
}
