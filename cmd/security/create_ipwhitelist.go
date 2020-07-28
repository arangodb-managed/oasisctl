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
		cmd.CreateCmd,
		&cobra.Command{
			Use:        "ipwhitelist",
			Short:      "Create a new IP whitelist",
			Deprecated: "Use ipallowlist instead",
			Hidden:     true,
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				name           string
				description    string
				organizationID string
				projectID      string
				cidrRanges     []string
			}{}
			f.StringVar(&cargs.name, "name", "", "Name of the IP whitelist")
			f.StringVar(&cargs.description, "description", "", "Description of the IP whitelist")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization to create the IP whitelist in")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project to create the IP whitelist in")
			f.StringSliceVar(&cargs.cidrRanges, "cidr-range", nil, "List of CIDR ranges from which deployments are accessible")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				name, argsUsed := cmd.ReqOption("name", cargs.name, args, 0)
				description := cargs.description
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				securityc := security.NewSecurityServiceClient(conn)
				rmc := rm.NewResourceManagerServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Fetch project
				project := selection.MustSelectProject(ctx, log, cargs.projectID, cargs.organizationID, rmc)

				// Create IP whitelist
				sort.Strings(cargs.cidrRanges)
				result, err := securityc.CreateIPWhitelist(ctx, &security.IPWhitelist{
					ProjectId:   project.GetId(),
					Name:        name,
					Description: description,
					CidrRanges:  cargs.cidrRanges,
				})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to create IP whitelist")
				}

				// Show result
				format.DisplaySuccess(cmd.RootArgs.Format)
				fmt.Println(format.IPWhitelist(result, cmd.RootArgs.Format))
			}
		},
	)
}
