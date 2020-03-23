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

package security

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
	security "github.com/arangodb-managed/apis/security/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.UpdateCmd,
		&cobra.Command{
			Use:   "ipwhitelist",
			Short: "Update an IP whitelist the authenticated user has access to",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				ipwhitelistID    string
				organizationID   string
				projectID        string
				name             string
				description      string
				addCidrRanges    []string
				removeCidrRanges []string
			}{}
			f.StringVarP(&cargs.ipwhitelistID, "ipwhitelist-id", "i", cmd.DefaultIPWhitelist(), "Identifier of the IP whitelist")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")
			f.StringVar(&cargs.name, "name", "", "Name of the CA certificate")
			f.StringVar(&cargs.description, "description", "", "Description of the CA certificate")
			f.StringSliceVar(&cargs.addCidrRanges, "add-cidr-range", nil, "List of CIDR ranges to add to the IP whitelist")
			f.StringSliceVar(&cargs.removeCidrRanges, "remove-cidr-range", nil, "List of CIDR ranges to remove from the IP whitelist")

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
				cidrRanges := make(map[string]struct{})
				for _, x := range item.GetCidrRanges() {
					cidrRanges[x] = struct{}{}
				}
				if len(cargs.addCidrRanges) > 0 {
					for _, x := range cargs.addCidrRanges {
						if _, found := cidrRanges[x]; !found {
							cidrRanges[x] = struct{}{}
							hasChanges = true
						}
					}
				}
				if len(cargs.removeCidrRanges) > 0 {
					for _, x := range cargs.removeCidrRanges {
						if _, found := cidrRanges[x]; found {
							delete(cidrRanges, x)
							hasChanges = true
						}
					}
				}
				if !hasChanges {
					fmt.Println("No changes")
				} else {
					// Rebuild CidrRanges list
					item.CidrRanges = make([]string, 0, len(cidrRanges))
					sort.Strings(item.CidrRanges)
					for x := range cidrRanges {
						item.CidrRanges = append(item.CidrRanges, x)
					}
					// Update IP whitelist
					updated, err := securityc.UpdateIPWhitelist(ctx, item)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to update IP whitelist")
					}

					// Show result
					fmt.Println("Updated IP whitelist!")
					fmt.Println(format.IPWhitelist(updated, cmd.RootArgs.Format))
				}
			}
		},
	)
}
