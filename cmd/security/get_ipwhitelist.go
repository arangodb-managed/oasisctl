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
		cmd.GetCmd,
		&cobra.Command{
			Use:   "ipwhitelist",
			Short: "Get an IP whitelist the authenticated user has access to",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				ipwhitelistID  string
				organizationID string
				projectID      string
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

				// Show result
				fmt.Println(format.IPWhitelist(item, cmd.RootArgs.Format))
			}
		},
	)
}
