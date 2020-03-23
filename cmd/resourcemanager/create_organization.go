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
)

var (
	// createOrganizationCmd creates a new organization
	createOrganizationCmd = &cobra.Command{
		Use:   "organization",
		Short: "Create a new organization",
		Run:   createOrganizationCmdRun,
	}
	createOrganizationArgs struct {
		name        string
		description string
	}
)

func init() {
	cmd.CreateCmd.AddCommand(createOrganizationCmd)

	f := createOrganizationCmd.Flags()
	f.StringVar(&createOrganizationArgs.name, "name", "", "Name of the organization")
	f.StringVar(&createOrganizationArgs.description, "description", "", "Description of the organization")
}

func createOrganizationCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := createOrganizationArgs
	name, argsUsed := cmd.ReqOption("name", cargs.name, args, 0)
	description := cargs.description
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Create organization
	result, err := rmc.CreateOrganization(ctx, &rm.Organization{
		Name:        name,
		Description: description,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create organization")
	}

	// Show result
	fmt.Println("Success!")
	fmt.Println(format.Organization(result, cmd.RootArgs.Format))
}
