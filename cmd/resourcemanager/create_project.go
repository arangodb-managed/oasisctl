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
	// createProjectCmd creates a new project
	createProjectCmd = &cobra.Command{
		Use:   "project",
		Short: "Create a new project",
		Run:   createProjectCmdRun,
	}
	createProjectArgs struct {
		name           string
		description    string
		organizationID string
	}
)

func init() {
	cmd.CreateCmd.AddCommand(createProjectCmd)

	f := createProjectCmd.Flags()
	f.StringVar(&createProjectArgs.name, "name", "", "Name of the project")
	f.StringVar(&createProjectArgs.description, "description", "", "Description of the project")
	f.StringVarP(&createProjectArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization to create the project in")
}

func createProjectCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := createProjectArgs
	name, argsUsed := cmd.ReqOption("name", cargs.name, args, 0)
	description := createProjectArgs.description
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch organization
	org := selection.MustSelectOrganization(ctx, log, cargs.organizationID, rmc)

	// Create project
	result, err := rmc.CreateProject(ctx, &rm.Project{
		OrganizationId: org.GetId(),
		Name:           name,
		Description:    description,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create project")
	}

	// Show result
	format.DisplaySuccess(cmd.RootArgs.Format)
	fmt.Println(format.Project(result, cmd.RootArgs.Format))
}
