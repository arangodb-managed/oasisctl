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
	// getOrganizationEmailCmd root for get organization email commands
	getOrganizationEmailCmd = &cobra.Command{
		Use:   "email",
		Short: "Get email specific information for an organization",
		Run:   cmd.ShowUsage,
	}

	// getOrganizationEmailDomainCmd root for get organization email domain commands
	getOrganizationEmailDomainCmd = &cobra.Command{
		Use:   "domain",
		Short: "Get email domain specific information for an organization",
		Run:   cmd.ShowUsage,
	}

	// getOrganizationDomainRestrictionsCmd fetches email domain restrictions that
	// are placed on an organization
	getOrganizationEmailDomainRestrictionsCmd = &cobra.Command{
		Use:   "restrictions",
		Short: "Get which email domain restrictions are placed on accessing a specific organization",
		Run:   getOrganizationEmailDomainRestrictionsCmdRun,
	}
	getOrganizationEmailDomainRestrictionsArgs struct {
		organizationID string
	}
)

func init() {
	getOrganizationCmd.AddCommand(getOrganizationEmailCmd)
	getOrganizationEmailCmd.AddCommand(getOrganizationEmailDomainCmd)
	getOrganizationEmailDomainCmd.AddCommand(getOrganizationEmailDomainRestrictionsCmd)
	f := getOrganizationEmailDomainRestrictionsCmd.Flags()
	f.StringVarP(&getOrganizationEmailDomainRestrictionsArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
}

func getOrganizationEmailDomainRestrictionsCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := getOrganizationEmailDomainRestrictionsArgs
	organizationID, argsUsed := cmd.OptOption("organization-id", cargs.organizationID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch organization
	item := selection.MustSelectOrganization(ctx, log, organizationID, rmc)

	// Show result
	fmt.Println(format.DomainRestrictions(item.GetEmailDomainRestrictions(), cmd.RootArgs.Format))
}
