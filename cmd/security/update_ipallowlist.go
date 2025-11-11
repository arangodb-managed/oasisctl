//
// DISCLAIMER
//
// Copyright 2020-2021 ArangoDB GmbH, Cologne, Germany
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
			Use:   "ipallowlist",
			Short: "Update an IP allowlist the authenticated user has access to",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				ipallowlistID           string
				organizationID          string
				projectID               string
				name                    string
				description             string
				addCidrRanges           []string
				removeCidrRanges        []string
				remoteInspectionAllowed bool
			}{}
			f.StringVarP(&cargs.ipallowlistID, "ipallowlist-id", "i", cmd.DefaultIPAllowlist(), "Identifier of the IP allowlist")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")
			f.StringVar(&cargs.name, "name", "", "Name of the CA certificate")
			f.StringVar(&cargs.description, "description", "", "Description of the CA certificate")
			f.StringSliceVar(&cargs.addCidrRanges, "add-cidr-range", nil, "List of CIDR ranges to add to the IP allowlist")
			f.StringSliceVar(&cargs.removeCidrRanges, "remove-cidr-range", nil, "List of CIDR ranges to remove from the IP allowlist")
			f.BoolVar(&cargs.remoteInspectionAllowed, "remote-inspection-allowed", false, "If set, remote connectivity checks by the Arango Managed Platform are allowed")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				ipallowlistID, argsUsed := cmd.OptOption("ipallowlist-id", cargs.ipallowlistID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				securityc := security.NewSecurityServiceClient(conn)
				rmc := rm.NewResourceManagerServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Fetch IP allowlist
				item := selection.MustSelectIPAllowlist(ctx, log, ipallowlistID, cargs.projectID, cargs.organizationID, securityc, rmc)

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
				existingCidrRanges := item.GetCidrRanges()

				for _, x := range existingCidrRanges {
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
				if f.Changed("remote-inspection-allowed") {
					item.RemoteInspectionAllowed = cargs.remoteInspectionAllowed
					hasChanges = true
				}
				if !hasChanges {
					fmt.Println("No changes")
				} else {
					// Rebuild CIDR ranges list
					UpdateCidrRanges(existingCidrRanges, cidrRanges, &UpdateCidrRangeOptions{
						AddedCidrRanges:   cargs.addCidrRanges,
						RemovedCidrRanges: cargs.removeCidrRanges,
					}, item)

					// Update IP allowlist
					updated, err := securityc.UpdateIPAllowlist(ctx, item)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to update IP allowlist")
					}

					// Show result
					fmt.Println("Updated IP allowlist!")
					fmt.Println(format.IPAllowlist(updated, cmd.RootArgs.Format))
				}
			}
		},
	)
}

// Update options for CIDR ranges
type UpdateCidrRangeOptions struct {
	AddedCidrRanges   []string
	RemovedCidrRanges []string
}

// UpdateCidrRanges updates the CIDR ranges to be added to the existing list without messing up the order of existing ones previously added
func UpdateCidrRanges(existingCidrRanges []string, cidrRanges map[string]struct{}, opts *UpdateCidrRangeOptions, item *security.IPAllowlist) {
	updatedCidrRanges := make([]string, 0, len(existingCidrRanges))
	for _, x := range existingCidrRanges {
		if _, found := cidrRanges[x]; found {
			updatedCidrRanges = append(updatedCidrRanges, x)
		}
	}

	if len(opts.AddedCidrRanges) > 0 {
		item.CidrRanges = append(updatedCidrRanges, opts.AddedCidrRanges...)
	} else {
		item.CidrRanges = updatedCidrRanges
	}

	if len(opts.RemovedCidrRanges) > 0 {
		updatedCidrRanges = make([]string, 0, len(item.CidrRanges))
		for _, x := range item.CidrRanges {
			if _, found := cidrRanges[x]; found {
				updatedCidrRanges = append(updatedCidrRanges, x)
			}
		}

		item.CidrRanges = updatedCidrRanges
	}
}
