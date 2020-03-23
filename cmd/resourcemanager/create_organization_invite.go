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

	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

var (
	// createOrganizationInvite creates a new organization invite
	createOrganizationInviteCmd = &cobra.Command{
		Use:   "invite",
		Short: "Create a new invite to an organization",
		Run:   createOrganizationInviteCmdRun,
	}
	createOrganizationInviteArgs struct {
		email          string
		organizationID string
	}
)

func init() {
	createOrganizationCmd.AddCommand(createOrganizationInviteCmd)

	f := createOrganizationInviteCmd.Flags()
	f.StringVar(&createOrganizationInviteArgs.email, "email", "", "Email address of the person to invite")
	f.StringVarP(&createOrganizationInviteArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization to create the invite in")
}

func createOrganizationInviteCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := createOrganizationInviteArgs
	email, argsUsed := cmd.ReqOption("email", cargs.email, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch organization
	org := selection.MustSelectOrganization(ctx, log, cargs.organizationID, rmc)

	// Create invite
	result, err := rmc.CreateOrganizationInvite(ctx, &rm.OrganizationInvite{
		OrganizationId: org.GetId(),
		Email:          email,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create organization invite")
	}

	// Show result
	fmt.Println("Success!")
	fmt.Println(format.OrganizationInvite(ctx, result, iamc, cmd.RootArgs.Format))
}
