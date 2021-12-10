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
	// updateOrganizationEmailCmd root for update organization email commands
	updateOrganizationEmailCmd = &cobra.Command{
		Use:   "email",
		Short: "Update email specific information for an organization",
		Run:   cmd.ShowUsage,
	}

	// updateOrganizationEmailDomainCmd root for update organization email domain commands
	updateOrganizationEmailDomainCmd = &cobra.Command{
		Use:   "domain",
		Short: "Update email domain specific information for an organization",
		Run:   cmd.ShowUsage,
	}

	// updateOrganizationEmailDomainRestrictionsCmd updates domain restrictions that
	// are placed on an organization
	updateOrganizationEmailDomainRestrictionsCmd = &cobra.Command{
		Use:   "restrictions",
		Short: "Update which domain restrictions are placed on accessing a specific organization",
		Run:   updateOrganizationEmailDomainRestrictionsCmdRun,
	}
	updateOrganizationEmailDomainRestrictionsArgs struct {
		organizationID string
		allowedDomains []string
	}
)

func init() {
	updateOrganizationCmd.AddCommand(updateOrganizationEmailCmd)
	updateOrganizationEmailCmd.AddCommand(updateOrganizationEmailDomainCmd)
	updateOrganizationEmailDomainCmd.AddCommand(updateOrganizationEmailDomainRestrictionsCmd)
	f := updateOrganizationEmailDomainRestrictionsCmd.Flags()
	f.StringVarP(&updateOrganizationEmailDomainRestrictionsArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
	f.StringSliceVarP(&updateOrganizationEmailDomainRestrictionsArgs.allowedDomains, "allowed-domain", "d", nil, "Allowed email domains for users of the organization")
}

func updateOrganizationEmailDomainRestrictionsCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := updateOrganizationEmailDomainRestrictionsArgs
	organizationID, argsUsed := cmd.OptOption("organization-id", cargs.organizationID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch organization
	item := selection.MustSelectOrganization(ctx, log, organizationID, rmc)

	// Set changes
	f := c.Flags()
	hasChanges := false
	if f.Changed("allowed-domain") {
		item.GetEmailDomainRestrictions().AllowedDomains = cargs.allowedDomains
		hasChanges = true
	}
	if !hasChanges {
		fmt.Println("No changes")
	} else {
		// Update project
		updated, err := rmc.UpdateOrganization(ctx, item)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to update organization email domain restrictions")
		}

		// Show result
		fmt.Println("Updated organization email domain restrictions!")
		fmt.Println(format.DomainRestrictions(updated.GetEmailDomainRestrictions(), cmd.RootArgs.Format))
	}
}
