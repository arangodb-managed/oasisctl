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
	// updateOrganizationAuthenticationProvidersCmd is base for "update organization authentication ..." commands
	updateOrganizationAuthenticationCmd = &cobra.Command{
		Use:   "authentication",
		Short: "Update authentication settings for an organization",
		Run:   cmd.ShowUsage,
	}
	// updateOrganizationAuthenticationProvidersCmd updates the allowed authentication providers
	// for an organization that the user has access to
	updateOrganizationAuthenticationProvidersCmd = &cobra.Command{
		Use:   "providers",
		Short: "Update allowed authentication providers for an organization the authenticated user has access to",
		Run:   updateOrganizationAuthenticationProvidersCmdRun,
	}
	updateOrganizationAuthenticationProvidersArgs struct {
		organizationID         string
		enableGithub           bool
		enableGoogle           bool
		enableMicrosoft        bool
		enableUsernamePassword bool
		enableSso              bool
	}
)

func init() {
	updateOrganizationCmd.AddCommand(updateOrganizationAuthenticationCmd)
	updateOrganizationAuthenticationCmd.AddCommand(updateOrganizationAuthenticationProvidersCmd)
	f := updateOrganizationAuthenticationProvidersCmd.Flags()
	f.StringVarP(&updateOrganizationAuthenticationProvidersArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
	f.BoolVar(&updateOrganizationAuthenticationProvidersArgs.enableGithub, "enable-github", false, "If set, allow access from user accounts authentication through Github")
	f.BoolVar(&updateOrganizationAuthenticationProvidersArgs.enableGoogle, "enable-google", false, "If set, allow access from user accounts authentication through Google")
	f.BoolVar(&updateOrganizationAuthenticationProvidersArgs.enableMicrosoft, "enable-microsoft", false, "If set, allow access from user accounts authentication through Microsoft")
	f.BoolVar(&updateOrganizationAuthenticationProvidersArgs.enableUsernamePassword, "enable-username-password", false, "If set, allow access from user accounts authentication through username-password")
	f.BoolVar(&updateOrganizationAuthenticationProvidersArgs.enableSso, "enable-sso", false, "If set, allow access from user accounts authentication through sso")
}

func updateOrganizationAuthenticationProvidersCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := updateOrganizationAuthenticationProvidersArgs
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
	if f.Changed("enable-github") {
		item.GetAuthenticationProviders().EnableGithub = cargs.enableGithub
		hasChanges = true
	}
	if f.Changed("enable-google") {
		item.GetAuthenticationProviders().EnableGoogle = cargs.enableGoogle
		hasChanges = true
	}
	if f.Changed("enable-microsoft") {
		item.GetAuthenticationProviders().EnableMicrosoft = cargs.enableMicrosoft
		hasChanges = true
	}
	if f.Changed("enable-username-password") {
		item.GetAuthenticationProviders().EnableUsernamePassword = cargs.enableUsernamePassword
		hasChanges = true
	}
	if f.Changed("enable-sso") {
		item.GetAuthenticationProviders().EnableSso = cargs.enableSso
		hasChanges = true
	}
	if !hasChanges {
		fmt.Println("No changes")
	} else {
		// Update project
		updated, err := rmc.UpdateOrganization(ctx, item)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to update organization authentication providers")
		}

		// Show result
		fmt.Println("Updated organization authentication providers!")
		fmt.Println(format.AuthenticationProviders(updated.GetAuthenticationProviders(), cmd.RootArgs.Format))
	}
}
