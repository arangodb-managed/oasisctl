//
// DISCLAIMER
//
// Copyright 2021 ArangoDB GmbH, Cologne, Germany
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

package network

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	common "github.com/arangodb-managed/apis/common/v1"
	data "github.com/arangodb-managed/apis/data/v1"
	network "github.com/arangodb-managed/apis/network/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

func init() {
	getPrivateEndpoint := &cobra.Command{
		Use:   "endpoint",
		Short: "",
		Run:   cmd.ShowUsage,
	}
	cmd.GetPrivateCmd.AddCommand(getPrivateEndpoint)

	cmd.InitCommand(
		getPrivateEndpoint,
		&cobra.Command{
			Use:   "service",
			Short: "Get a Private Endpoint Service the authenticated user has access to",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				deplID         string
				organizationID string
				projectID      string
			}{}
			f.StringVarP(&cargs.deplID, "deployment-id", "d", cmd.DefaultDeployment(), "Identifier of the deployment that the private endpoint service is connected to")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				deplID, argsUsed := cmd.OptOption("deployment-id", cargs.deplID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				datac := data.NewDataServiceClient(conn)
				nwc := network.NewNetworkServiceClient(conn)
				rmc := rm.NewResourceManagerServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Fetch deployment
				depl := selection.MustSelectDeployment(ctx, log, deplID, cargs.projectID, cargs.organizationID, datac, rmc)

				// Fetch private endpoint service
				item, err := nwc.GetPrivateEndpointServiceByDeploymentID(ctx, &common.IDOptions{Id: depl.GetId()})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to get private endpoint service")
				}
				// Show result
				fmt.Println(format.PrivateEndpointService(item, cmd.RootArgs.Format))
			}

		},
	)
}
