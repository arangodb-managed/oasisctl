//
// DISCLAIMER
//
// Copyright 2021-2022 ArangoDB GmbH, Cologne, Germany
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

package network

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	common "github.com/arangodb-managed/apis/common/v1"
	data "github.com/arangodb-managed/apis/data/v1"
	network "github.com/arangodb-managed/apis/network/v1"
	platform "github.com/arangodb-managed/apis/platform/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

func init() {
	updatePrivateEndpoint := &cobra.Command{
		Use:   "endpoint",
		Short: "",
		Run:   cmd.ShowUsage,
	}
	cmd.UpdatePrivateCmd.AddCommand(updatePrivateEndpoint)

	cmd.InitCommand(
		updatePrivateEndpoint,
		&cobra.Command{
			Use:   "service",
			Short: "Update a Private Endpoint Service attached to an existing deployment",
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
				awsPrincipals           []string
				gcpProjects             []string
			}{}
			f.StringVarP(&cargs.deplID, "deployment-id", "d", cmd.DefaultDeployment(), "Identifier of the deployment that the private endpoint service is connected to")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")
			f.StringVar(&cargs.name, "name", "", "Name of the private endpoint service")
			f.StringVar(&cargs.description, "description", "", "Description of the private endpoint service")
			f.StringSliceVar(&cargs.alternateDNSNames, "alternate-dns-name", nil, "DNS names used for the deployment in the private network")
			f.StringSliceVar(&cargs.azClientSubscriptionIDs, "azure-client-subscription-id", nil, "List of Azure subscription IDs from which a Private Endpoint can be created")
			f.StringSliceVar(&cargs.awsPrincipals, "aws-principal", nil, "List of AWS Principals from which a Private Endpoint can be created (Format: <AccountID>[/Role/<RoleName>|/User/<UserName>])")
			f.StringSliceVar(&cargs.gcpProjects, "google-project", nil, "List of Google projects from which a Private Endpoint can be created")

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
				platformc := platform.NewPlatformServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Fetch deployment
				depl := selection.MustSelectDeployment(ctx, log, deplID, cargs.projectID, cargs.organizationID, datac, rmc)

				// Fetch region
				region, err := platformc.GetRegion(ctx, &common.IDOptions{Id: depl.GetRegionId()})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to get region for deployment")
				}

				// Fetch existing private endpoint service
				item, err := nwc.GetPrivateEndpointServiceByDeploymentID(ctx, &common.IDOptions{Id: depl.GetId()})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to get private endpoint service")
				}

				// Set changes
				f := c.Flags()
				hasChanges := false
				if f.Changed("name") {
					item.Name = cargs.name
					hasChanges = true
				}
				if f.Changed("description") {
					item.Description = cargs.description
					hasChanges = true
				}
				if f.Changed("alternate-dns-name") {
					item.AlternateDnsNames = cargs.alternateDNSNames
					hasChanges = true
				}
				switch region.GetProviderId() {
				case "aks":
					if f.Changed("azure-client-subscription-id") {
						if item.Aks == nil {
							item.Aks = &network.PrivateEndpointService_Aks{}
						}
						item.Aks.ClientSubscriptionIds = cargs.azClientSubscriptionIDs
						hasChanges = true
					}
				case "aws":
					if f.Changed("aws-principal") {
						if item.Aws == nil {
							item.Aws = &network.PrivateEndpointService_Aws{}
						}
						p, err := getAwsPrincipals(cargs.awsPrincipals)
						if err != nil {
							log.Fatal().Err(err).Msg("Failed to parse AWS principals")
						}
						item.Aws.AwsPrincipals = p
						hasChanges = true
					}
				case "gcp":
					if f.Changed("google-project") {
						if item.Gcp == nil {
							item.Gcp = &network.PrivateEndpointService_Gcp{}
						}
						item.Gcp.Projects = cargs.gcpProjects
						hasChanges = true
					}
				}

				if !hasChanges {
					fmt.Println("No changes")
				} else {
					// Update private endpoint service
					_, err = nwc.UpdatePrivateEndpointService(ctx, item)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to update private endpoint service")
					}

					// Show result
					fmt.Println("Updated private endpoint service")
				}
			}

		},
	)
}
