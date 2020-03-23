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
// Author Ewout Prangsma
//

package platform

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	platform "github.com/arangodb-managed/apis/platform/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.GetCmd,
		&cobra.Command{
			Use:   "region",
			Short: "Get a region the authenticated user has access to",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				organizationID string
				regionID       string
				providerID     string
			}{}
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Optional Identifier of the organization")
			f.StringVarP(&cargs.regionID, "region-id", "r", cmd.DefaultRegion(), "Identifier of the region")
			f.StringVarP(&cargs.providerID, "provider-id", "p", cmd.DefaultProvider(), "Identifier of the provider")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				regionID, argsUsed := cmd.OptOption("region-id", cargs.regionID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				platformc := platform.NewPlatformServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Fetch region
				item := selection.MustSelectRegion(ctx, log, regionID, cargs.providerID, cargs.organizationID, platformc)

				// Show result
				fmt.Println(format.Region(item, cmd.RootArgs.Format))
			}
		},
	)
}
