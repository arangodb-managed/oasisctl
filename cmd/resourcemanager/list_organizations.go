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
	"github.com/arangodb-managed/oasisctl/pkg/format"
)

var (
	// listOrganizationsCmd fetches organizations the user is a part of
	listOrganizationsCmd = &cobra.Command{
		Use:   "organizations",
		Short: "List all organizations the authenticated user is a member of",
		Run:   listOrganizationsCmdRun,
	}
)

func init() {
	cmd.ListCmd.AddCommand(listOrganizationsCmd)
}

func listOrganizationsCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cmd.MustCheckNumberOfArgs(args, 0)

	// Connect
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch organizations
	list, err := rmc.ListOrganizations(ctx, &common.ListOptions{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list organizations")
	}

	// Show result
	fmt.Println(format.OrganizationList(list.Items, cmd.RootArgs.Format))
}
