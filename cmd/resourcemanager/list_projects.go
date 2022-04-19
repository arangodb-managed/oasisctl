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
	"context"
	"fmt"

	"github.com/spf13/cobra"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

var (
	// listProjectsCmd fetches projects of the given organization
	listProjectsCmd = &cobra.Command{
		Use:   "projects",
		Short: "List all projects of the given organization",
		Run:   listProjectsCmdRun,
	}
	listProjectsArgs struct {
		organizationID string
	}
)

func init() {
	cmd.ListCmd.AddCommand(listProjectsCmd)
	f := listProjectsCmd.Flags()
	f.StringVarP(&listProjectsArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
}

func listProjectsCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := listProjectsArgs
	organizationID, argsUsed := cmd.OptOption("organization-id", cargs.organizationID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch organization
	org := selection.MustSelectOrganization(ctx, log, organizationID, rmc)

	// Fetch projects in organizationo
	var projects []*rm.Project
	if err := rm.ForEachProject(ctx, func(ctx context.Context, req *common.ListOptions) (*rm.ProjectList, error) {
		return rmc.ListProjects(ctx, req)
	}, &common.ListOptions{ContextId: org.GetId()}, func(ctx context.Context, project *rm.Project) error {
		projects = append(projects, project)
		return nil
	}); err != nil {
		log.Fatal().Err(err).Msg("Failed to list projects")
	}

	// Show result
	fmt.Println(format.ProjectList(projects, cmd.RootArgs.Format))
}
