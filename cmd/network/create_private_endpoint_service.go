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
	"strings"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	common "github.com/arangodb-managed/apis/common/v1"
	data "github.com/arangodb-managed/apis/data/v1"
	network "github.com/arangodb-managed/apis/network/v1"
	platform "github.com/arangodb-managed/apis/platform/v1"
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
			f.StringSliceVar(&cargs.gcpProjects, "gcp-project", nil, "List of GCP projects from which a Private Endpoint can be created")

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

				// Generate default arguments
				if cargs.name == "" {
					cargs.name = "Private Endpoint Service for " + depl.GetName()
				}

				// Make object
				pes := &network.PrivateEndpointService{
					DeploymentId:      depl.GetId(),
					Name:              cargs.name,
					Description:       cargs.description,
					AlternateDnsNames: cargs.alternateDNSNames,
				}

				switch region.GetProviderId() {
				case "aks":
					pes.Aks = &network.PrivateEndpointService_Aks{
						ClientSubscriptionIds: cargs.azClientSubscriptionIDs,
					}
				case "aws":
					p, err := getAwsPrincipals(cargs.awsPrincipals)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to parse AWS principals")
					}
					pes.Aws = &network.PrivateEndpointService_Aws{
						AwsPrincipals: p,
					}
				case "gcp":
					pes.Gcp = &network.PrivateEndpointService_Gcp{
						Projects: cargs.gcpProjects,
					}
				}

				// Create private endpoint service
				item, err := nwc.CreatePrivateEndpointService(ctx, pes)
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to create private endpoint service")
				}
				// Show result
				fmt.Println(format.PrivateEndpointService(item, cmd.RootArgs.Format))
			}

		},
	)
}

// getAwsPrincipals get the AWS principals out of a string slice
// with format: <AccountID>[/Role/<RoleName>|/User/<UserName>]
func getAwsPrincipals(source []string) ([]*network.PrivateEndpointService_AwsPrincipals, error) {
	principals := make(map[string]*network.PrivateEndpointService_AwsPrincipals)
	for _, s := range source {
		p, err := parseAwsPrincipal(s)
		if err != nil {
			return nil, err
		}
		principal, found := principals[p.AccountID]
		if !found {
			principal = &network.PrivateEndpointService_AwsPrincipals{
				AccountId: p.AccountID,
			}
			principals[p.AccountID] = principal
		}
		if p.RoleName != "" {
			principal.RoleNames = append(principal.RoleNames, p.RoleName)
		}
		if p.UserName != "" {
			principal.UserNames = append(principal.UserNames, p.UserName)
		}
	}
	// Convert map into slice
	var result []*network.PrivateEndpointService_AwsPrincipals
	for _, p := range principals {
		result = append(result, p)
	}
	return result, nil
}

type parsedAwsPrincipal struct {
	// Required account ID (12 digits)
	AccountID string
	// Optional role name
	RoleName string
	// Optional user name
	UserName string
}

// parseAwsPrincipal convert a string into strong typed object.
// input in format: <AccountID>[/Role/<RoleName>|/User/<UserName>]
func parseAwsPrincipal(input string) (parsedAwsPrincipal, error) {
	input = strings.TrimSpace(input)
	splitted := strings.Split(input, "/")
	if len(splitted) == 1 {
		return parsedAwsPrincipal{
			AccountID: splitted[0],
		}, nil
	}
	if len(splitted) != 3 {
		return parsedAwsPrincipal{}, fmt.Errorf("cannot parse AWS principal: %s", input)
	}
	switch splitted[1] {
	case "Role":
		return parsedAwsPrincipal{
			AccountID: splitted[0],
			RoleName:  splitted[2],
		}, nil
	case "User":
		return parsedAwsPrincipal{
			AccountID: splitted[0],
			UserName:  splitted[2],
		}, nil
	default:
		return parsedAwsPrincipal{}, fmt.Errorf("cannot parse AWS principal: %s (unknown subtype)", input)
	}
}
