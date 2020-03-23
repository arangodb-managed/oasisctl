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
	// updateOrganizationCmd updates an organization that the user has access to
	updateOrganizationCmd = &cobra.Command{
		Use:   "organization",
		Short: "Update an organization the authenticated user has access to",
		Run:   updateOrganizationCmdRun,
	}
	updateOrganizationArgs struct {
		organizationID string
		name           string
		description    string
	}
)

func init() {
	cmd.UpdateCmd.AddCommand(updateOrganizationCmd)
	f := updateOrganizationCmd.Flags()
	f.StringVarP(&updateOrganizationArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
	f.StringVar(&updateOrganizationArgs.name, "name", "", "Name of the organization")
	f.StringVar(&updateOrganizationArgs.description, "description", "", "Description of the organization")
}

func updateOrganizationCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := updateOrganizationArgs
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
	if f.Changed("name") {
		item.Name = cargs.name
		hasChanges = true
	}
	if f.Changed("description") {
		item.Description = cargs.description
		hasChanges = true
	}
	if !hasChanges {
		fmt.Println("No changes")
	} else {
		// Update project
		updated, err := rmc.UpdateOrganization(ctx, item)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to update organization")
		}

		// Show result
		fmt.Println("Updated organization!")
		fmt.Println(format.Organization(updated, cmd.RootArgs.Format))
	}
}
