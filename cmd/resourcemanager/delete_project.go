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

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

var (
	// deleteProjectCmd deletes a project that the user has access to
	deleteProjectCmd = &cobra.Command{
		Use:   "project",
		Short: "Delete a project the authenticated user has access to",
		Run:   deleteProjectCmdRun,
	}
	deleteProjectArgs struct {
		organizationID string
		projectID      string
	}
)

func init() {
	cmd.DeleteCmd.AddCommand(deleteProjectCmd)
	f := deleteProjectCmd.Flags()
	f.StringVarP(&deleteProjectArgs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")
	f.StringVarP(&deleteProjectArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
}

func deleteProjectCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := deleteProjectArgs
	projectID, argsUsed := cmd.OptOption("project-id", cargs.projectID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch project
	item := selection.MustSelectProject(ctx, log, projectID, cargs.organizationID, rmc)

	// Delete project
	if _, err := rmc.DeleteProject(ctx, &common.IDOptions{Id: item.GetId()}); err != nil {
		log.Fatal().Err(err).Msg("Failed to delete project")
	}

	// Show result
	fmt.Println("Deleted project!")
}
