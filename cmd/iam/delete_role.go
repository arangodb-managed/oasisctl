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

	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

var (
	// deleteRoleCmd deletes a role that the user has access to
	deleteRoleCmd = &cobra.Command{
		Use:   "role",
		Short: "Delete a role the authenticated user has access to",
		Run:   deleteRoleCmdRun,
	}
	deleteRoleArgs struct {
		organizationID string
		roleID         string
	}
)

func init() {
	cmd.DeleteCmd.AddCommand(deleteRoleCmd)
	f := deleteRoleCmd.Flags()
	f.StringVarP(&deleteRoleArgs.roleID, "role-id", "r", cmd.DefaultRole(), "Identifier of the role")
	f.StringVarP(&deleteRoleArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
}

func deleteRoleCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := deleteRoleArgs
	roleID, argsUsed := cmd.OptOption("role-id", cargs.roleID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch role
	item := selection.MustSelectRole(ctx, log, roleID, cargs.organizationID, iamc, rmc)

	// Delete role
	if _, err := iamc.DeleteRole(ctx, &common.IDOptions{Id: item.GetId()}); err != nil {
		log.Fatal().Err(err).Msg("Failed to delete role")
	}

	// Show result
	fmt.Println("Deleted role!")
}
