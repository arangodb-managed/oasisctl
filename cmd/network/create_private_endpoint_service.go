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

	data "github.com/arangodb-managed/apis/data/v1"
	network "github.com/arangodb-managed/apis/network/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

func init() {
	createPrivateEndpoint := &cobra.Command{
		Use:   "endpoint",
		Short: "",
		Run:   cmd.ShowUsage,
	}
	cmd.CreatePrivateCmd.AddCommand(createPrivateEndpoint)

	cmd.InitCommand(
		createPrivateEndpoint,
		&cobra.Command{
			Use:   "service",
			Short: "Create a Private Endpoint Service attached to an existing deployment",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				deplID                  string
				organizationID          string
				projectID               string
				name                    string
				description             string
				alternateDNSNames       []string
				azClientSubscriptionIDs []string
			}{}
			f.StringVarP(&cargs.deplID, "deployment-id", "d", cmd.DefaultDeployment(), "Identifier of the deployment that the private endpoint service is connected to")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")
			f.StringVar(&cargs.name, "name", "", "Name of the private endpoint service")
			f.StringVar(&cargs.description, "description", "", "Description of the private endpoint service")
			f.StringSliceVar(&cargs.alternateDNSNames, "alternate-dns-name", nil, "DNS names used for the deployment in the private network")
			f.StringSliceVar(&cargs.azClientSubscriptionIDs, "azure-client-subscription-id", nil, "List of Azure subscription IDs from which a Private Endpoint can be created")

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

				// Generate default arguments
				if cargs.name == "" {
					cargs.name = "Private Endpoint Service for " + depl.GetName()
				}

				// Create private endpoint service
				item, err := nwc.CreatePrivateEndpointService(ctx, &network.PrivateEndpointService{
					DeploymentId:      depl.GetId(),
					Name:              cargs.name,
					Description:       cargs.description,
					AlternateDnsNames: cargs.alternateDNSNames,
					Aks: &network.PrivateEndpointService_Aks{
						ClientSubscriptionIds: cargs.azClientSubscriptionIDs,
					},
				})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to create private endpoint service")
				}
				// Show result
				fmt.Println(format.PrivateEndpointService(item, cmd.RootArgs.Format))
			}

		},
	)
}
