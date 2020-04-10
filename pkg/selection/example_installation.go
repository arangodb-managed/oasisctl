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

package selection

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	common "github.com/arangodb-managed/apis/common/v1"
	data "github.com/arangodb-managed/apis/data/v1"
	example "github.com/arangodb-managed/apis/example/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

// MustSelectExampleDatasetInstallation fetches an example dataset installation with given ID, name, or URL and fails if no such item is found.
// If no ID is specified, all installations are fetched from the selected deployment
// and if the list is exactly 1 long, that backup is returned.
func MustSelectExampleDatasetInstallation(ctx context.Context, log zerolog.Logger, id, deploymentID, projectID, organizationID string,
	datac data.DataServiceClient,
	examplec example.ExampleDatasetServiceClient,
	rmc rm.ResourceManagerServiceClient) *example.ExampleDatasetInstallation {
	result, err := SelectExampleDatasetInstallation(ctx, log, id, deploymentID, projectID, organizationID, datac, examplec, rmc)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to select example dataset installation")
	}
	return result
}

// SelectExampleDatasetInstallation fetches an example dataset with given ID, name, or URL or returns an error if not found.
func SelectExampleDatasetInstallation(ctx context.Context, log zerolog.Logger, id, deploymentID, projectID, organizationID string,
	datac data.DataServiceClient,
	examplec example.ExampleDatasetServiceClient,
	rmc rm.ResourceManagerServiceClient) (*example.ExampleDatasetInstallation, error) {
	if id == "" {
		depl, err := SelectDeployment(ctx, log, deploymentID, projectID, organizationID, datac, rmc)
		if err != nil {
			return nil, err
		}
		list, err := examplec.ListExampleDatasetInstallations(ctx, &example.ListExampleDatasetInstallationsRequest{DeploymentId: depl.GetId()})
		if err != nil {
			log.Debug().Err(err).Msg("Failed to list example dataset installations")
			return nil, err
		}
		if len(list.Items) != 1 {
			log.Debug().Err(err).Msgf("You have access to %d example dataset installations. Please specify one explicitly.", len(list.Items))
			return nil, fmt.Errorf("You have access to %d example dataset installations. Please specify one explicitly.", len(list.Items))
		}
		return list.Items[0], nil
	}
	result, err := examplec.GetExampleDatasetInstallation(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) {
			// Try to lookup example dataset installation by URL
			depl, err := SelectDeployment(ctx, log, deploymentID, projectID, organizationID, datac, rmc)
			if err != nil {
				return nil, err
			}
			list, err := examplec.ListExampleDatasetInstallations(ctx, &example.ListExampleDatasetInstallationsRequest{DeploymentId: depl.GetId()})
			if err == nil {
				for _, x := range list.Items {
					if x.GetUrl() == id {
						return x, nil
					}
				}
			}
		}
		log.Debug().Err(err).Str("example_dataset_installation", id).Msg("Failed to get example dataset installation")
		return nil, err
	}
	return result, nil
}
