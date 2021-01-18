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
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

// MustSelectDeployment fetches the deployment given ID, name, or URL and fails if no deployment is found.
// If no ID is specified, all deployments are fetched from the selected project
// and if the list is exactly 1 long, that deployment is returned.
func MustSelectDeployment(ctx context.Context, log zerolog.Logger, id, projectID, orgID string, datac data.DataServiceClient, rmc rm.ResourceManagerServiceClient) *data.Deployment {
	deployment, err := SelectDeployment(ctx, log, id, projectID, orgID, datac, rmc)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list deployments")
	}
	return deployment
}

// SelectDeployment fetches the deployment given ID, name, or URL or returns an error if not found.
// If no ID is specified, all deployments are fetched from the selected project
// and if the list is exactly 1 long, that deployment is returned.
func SelectDeployment(ctx context.Context, log zerolog.Logger, id, projectID, orgID string, datac data.DataServiceClient, rmc rm.ResourceManagerServiceClient) (*data.Deployment, error) {
	if id == "" {
		project, err := SelectProject(ctx, log, projectID, orgID, rmc)
		if err != nil {
			return nil, err
		}
		list, err := datac.ListDeployments(ctx, &common.ListOptions{ContextId: project.GetId()})
		if err != nil {
			log.Debug().Err(err).Msg("Failed to list deployments")
			return nil, err
		}
		if len(list.Items) != 1 {
			log.Debug().Err(err).Msgf("You have access to %d deployments. Please specify one explicitly.", len(list.Items))
			return nil, fmt.Errorf("You have access to %d deployments. Please specify one explicitly.", len(list.Items))

		}
		return list.Items[0], nil
	}
	result, err := datac.GetDeployment(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) {
			log.Debug().Msg("Deployment not found")
			// Try to lookup deployment by name or URL
			project, err := SelectProject(ctx, log, projectID, orgID, rmc)
			if err != nil {
				return nil, err
			}
			list, err := datac.ListDeployments(ctx, &common.ListOptions{ContextId: project.GetId()})
			if err == nil {
				for _, x := range list.Items {
					if x.GetName() == id || x.GetUrl() == id {
						return x, nil
					}
				}
			}
		}
		log.Debug().Err(err).Str("deployment", id).Msg("Failed to get deployment")
		return nil, err
	}
	return result, nil
}
