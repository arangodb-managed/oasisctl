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

package iam

import (
	"fmt"

	"github.com/spf13/cobra"

	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

var (
	// updateGroupCmd updates a group that the user has access to
	updateGroupCmd = &cobra.Command{
		Use:   "group",
		Short: "Update a group the authenticated user has access to",
		Run:   updateGroupCmdRun,
	}
	updateGroupArgs struct {
		groupID        string
		organizationID string
		name           string
		description    string
	}
)

func init() {
	cmd.UpdateCmd.AddCommand(updateGroupCmd)
	f := updateGroupCmd.Flags()
	f.StringVarP(&updateGroupArgs.groupID, "group-id", "g", cmd.DefaultGroup(), "Identifier of the group")
	f.StringVarP(&updateGroupArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
	f.StringVar(&updateGroupArgs.name, "name", "", "Name of the group")
	f.StringVar(&updateGroupArgs.description, "description", "", "Description of the group")
}

func updateGroupCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := updateGroupArgs
	groupID, argsUsed := cmd.OptOption("group-id", cargs.groupID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch group
	item := selection.MustSelectGroup(ctx, log, groupID, cargs.organizationID, iamc, rmc)

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
		// Update group
		updated, err := iamc.UpdateGroup(ctx, item)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to update group")
		}

		// Show result
		fmt.Println("Updated group!")
		fmt.Println(format.Group(updated, cmd.RootArgs.Format))
	}
}
