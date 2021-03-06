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
		cmd.UnlockCmd,
		&cobra.Command{
			Use:   "ipallowlist",
			Short: "Unlock an IP allowlist, so it can be deleted",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				ipallowlistID  string
				organizationID string
				projectID      string
			}{}
			f.StringVarP(&cargs.ipallowlistID, "ipallowlist-id", "i", cmd.DefaultIPAllowlist(), "Identifier of the IP allowlist")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")

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
				item.Locked = false
				// Update IP allowlist
				updated, err := securityc.UpdateIPAllowlist(ctx, item)
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to unlock IP allowlist")
				}

				// Show result
				fmt.Println("Unlocked IP allowlist!")
				fmt.Println(format.IPAllowlist(updated, cmd.RootArgs.Format))
			}
		},
	)
}
