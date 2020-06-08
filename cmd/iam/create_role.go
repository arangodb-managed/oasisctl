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
	// createRoleCmd creates a new role
	createRoleCmd = &cobra.Command{
		Use:   "role",
		Short: "Create a new role",
		Run:   createRoleCmdRun,
	}
	createRoleArgs struct {
		name           string
		description    string
		organizationID string
		permissions    []string
	}
)

func init() {
	cmd.CreateCmd.AddCommand(createRoleCmd)

	f := createRoleCmd.Flags()
	f.StringVar(&createRoleArgs.name, "name", "", "Name of the role")
	f.StringVar(&createRoleArgs.description, "description", "", "Description of the role")
	f.StringSliceVarP(&createRoleArgs.permissions, "permission", "p", nil, "Permissions granted by the role")
	f.StringVarP(&createRoleArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization to create the role in")
}

func createRoleCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := createRoleArgs
	name, argsUsed := cmd.ReqOption("name", cargs.name, args, 0)
	description := cargs.description
	permissions := cargs.permissions
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch organization
	org := selection.MustSelectOrganization(ctx, log, cargs.organizationID, rmc)

	// Create role
	result, err := iamc.CreateRole(ctx, &iam.Role{
		OrganizationId: org.GetId(),
		Name:           name,
		Description:    description,
		Permissions:    permissions,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create role")
	}

	// Show result
	format.DisplaySuccess(cmd.RootArgs.Format)
	fmt.Println(format.Role(result, cmd.RootArgs.Format))
}
