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
	metrics "github.com/arangodb-managed/apis/metrics/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

// MustSelectMetricsToken fetches the metrics token given ID, name, or URL and fails if no token is found.
// If no ID is specified, all tokens are fetched from the selected deployment
// and if the list is exactly 1 long, that token is returned.
func MustSelectMetricsToken(ctx context.Context, log zerolog.Logger, id, deploymentID, projectID, orgID string, metricsc metrics.MetricsServiceClient, datac data.DataServiceClient, rmc rm.ResourceManagerServiceClient) *metrics.Token {
	token, err := SelectMetricsToken(ctx, log, id, deploymentID, projectID, orgID, metricsc, datac, rmc)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list metrics tokens")
	}
	return token
}

// SelectMetricsToken fetches the metrics token given ID, name, or URL or returns an error if not found.
// If no ID is specified, all tokens are fetched from the selected deployment
// and if the list is exactly 1 long, that token is returned.
func SelectMetricsToken(ctx context.Context, log zerolog.Logger, id, deploymentID, projectID, orgID string, metricsc metrics.MetricsServiceClient, datac data.DataServiceClient, rmc rm.ResourceManagerServiceClient) (*metrics.Token, error) {
	if id == "" {
		deployment, err := SelectDeployment(ctx, log, deploymentID, projectID, orgID, datac, rmc)
		if err != nil {
			return nil, err
		}
		list, err := metricsc.ListTokens(ctx, &metrics.ListTokensRequest{DeploymentId: deployment.GetId()})
		if err != nil {
			log.Debug().Err(err).Msg("Failed to list metrics tokens")
			return nil, err
		}
		if len(list.Items) != 1 {
			log.Debug().Err(err).Msgf("You have access to %d metrics tokens. Please specify one explicitly.", len(list.Items))
			return nil, fmt.Errorf("You have access to %d metrics tokens. Please specify one explicitly.", len(list.Items))

		}
		return list.Items[0], nil
	}
	result, err := metricsc.GetToken(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) {
			// Try to lookup token by name or URL
			deployment, err := SelectDeployment(ctx, log, deploymentID, projectID, orgID, datac, rmc)
			if err != nil {
				return nil, err
			}
			list, err := metricsc.ListTokens(ctx, &metrics.ListTokensRequest{DeploymentId: deployment.GetId()})
			if err == nil {
				for _, x := range list.Items {
					if x.GetName() == id || x.GetUrl() == id {
						return x, nil
					}
				}
			}
		}
		log.Debug().Err(err).Str("token", id).Msg("Failed to get metrics token")
		return nil, err
	}
	return result, nil
}
