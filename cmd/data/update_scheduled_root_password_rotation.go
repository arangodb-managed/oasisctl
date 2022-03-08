//
// DISCLAIMER
//
// Copyright 2022 ArangoDB GmbH, Cologne, Germany
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

package data

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	data "github.com/arangodb-managed/apis/data/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.EnableCmd,
		&cobra.Command{
			Use:   "scheduled-root-password-rotation",
			Short: "Enable scheduled root password rotation",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			runFunc(c, f, true)
		},
	)

	cmd.InitCommand(
		cmd.DisableCmd,
		&cobra.Command{
			Use:   "scheduled-root-password-rotation",
			Short: "Disable scheduled root password rotation",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			runFunc(c, f, false)
		},
	)
}

func runFunc(c *cobra.Command, f *flag.FlagSet, enabled bool) {
	cargs := &struct {
		deploymentID   string
		organizationID string
		projectID      string
	}{}
	f.StringVarP(&cargs.deploymentID, "deployment-id", "d", cmd.DefaultDeployment(), "Identifier of the deployment")
	f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
	f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")

	c.Run = func(c *cobra.Command, args []string) {
		// Validate arguments
		log := cmd.CLILog
		deploymentID, argsUsed := cmd.OptOption("deployment-id", cargs.deploymentID, args, 0)
		cmd.MustCheckNumberOfArgs(args, argsUsed)

		// Connect
		conn := cmd.MustDialAPI()
		datac := data.NewDataServiceClient(conn)
		rmc := rm.NewResourceManagerServiceClient(conn)
		ctx := cmd.ContextWithToken()

		// Fetch deployment
		item := selection.MustSelectDeployment(ctx, log, deploymentID, cargs.projectID, cargs.organizationID, datac, rmc)

		action := "enable"
		if !enabled {
			action = "disable"
		}

		// Update
		_, err := datac.UpdateDeploymentScheduledRootPasswordRotation(ctx, &data.UpdateDeploymentScheduledRootPasswordRotationRequest{
			DeploymentId: item.GetId(),
			Enabled:      enabled,
		})
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to %s scheduled root password rotation.", action)
		}
		fmt.Printf("Scheduled root password rotation %sd.\n", action)
	}
}
