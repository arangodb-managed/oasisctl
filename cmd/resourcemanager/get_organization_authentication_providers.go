//
// DISCLAIMER
//
// Copyright 2021 ArangoDB GmbH, Cologne, Germany
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

package rm

import (
	"fmt"

	"github.com/spf13/cobra"

	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

var (
	// getOrganizationAuthenticationCmd root for get organization authentication commands
	getOrganizationAuthenticationCmd = &cobra.Command{
		Use:   "authentication",
		Short: "Get authentication specific information for an organization",
		Run:   cmd.ShowUsage,
	}

	// getOrganizationAuthenticationProvidersCmd fetches authentication providers that
	// are allowed for an organization
	getOrganizationAuthenticationProvidersCmd = &cobra.Command{
		Use:   "providers",
		Short: "Get which authentication providers are allowed for accessing a specific organization",
		Run:   getOrganizationAuthenticationProvidersCmdRun,
	}
	getOrganizationAuthenticationProvidersArgs struct {
		organizationID string
	}
)

func init() {
	getOrganizationCmd.AddCommand(getOrganizationAuthenticationCmd)
	getOrganizationAuthenticationCmd.AddCommand(getOrganizationAuthenticationProvidersCmd)
	f := getOrganizationAuthenticationProvidersCmd.Flags()
	f.StringVarP(&getOrganizationAuthenticationProvidersArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
}

func getOrganizationAuthenticationProvidersCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := getOrganizationAuthenticationProvidersArgs
	organizationID, argsUsed := cmd.OptOption("organization-id", cargs.organizationID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch organization
	item := selection.MustSelectOrganization(ctx, log, organizationID, rmc)

	// Show result
	fmt.Println(format.AuthenticationProviders(item.GetAuthenticationProviders(), cmd.RootArgs.Format))
}
