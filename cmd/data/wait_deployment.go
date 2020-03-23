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

package data

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	data "github.com/arangodb-managed/apis/data/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

const (
	defaultWaitDeploymentTimeout = time.Minute * 20
)

func init() {
	cmd.InitCommand(
		cmd.WaitCmd,
		&cobra.Command{
			Use:   "deployment",
			Short: "Wait for a deployment to reach the ready status",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				deploymentID   string
				organizationID string
				projectID      string
				timeout        time.Duration
			}{}
			f.StringVarP(&cargs.deploymentID, "deployment-id", "d", cmd.DefaultDeployment(), "Identifier of the deployment")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")
			f.DurationVarP(&cargs.timeout, "timeout", "t", defaultWaitDeploymentTimeout, "How long to wait for the deployment to reach the ready status")

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

				start := time.Now()
				for {
					// Fetch deployment
					item := selection.MustSelectDeployment(ctx, log, deploymentID, cargs.projectID, cargs.organizationID, datac, rmc)

					// Check status
					status := item.GetStatus()
					if status.GetReady() {
						// Status ready
						break
					}

					// Check timeout
					if time.Since(start) > cargs.timeout {
						log.Fatal().Msg("Deployment not ready after timeout")
					}

					// Wait a bit
					log.Debug().Str("status", status.GetDescription()).Msg("Deployment not yet ready")
					time.Sleep(time.Second * 2)
				}

				fmt.Println("Deployment ready")
			}
		},
	)
}
