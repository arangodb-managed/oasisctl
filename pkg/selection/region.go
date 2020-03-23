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
	platform "github.com/arangodb-managed/apis/platform/v1"
)

// MustSelectRegion fetches the region with given ID or location and fails if no region is found.
// If no ID is specified, all regions are fetched from the selected provider
// and if the list is exactly 1 long, that region is returned.
func MustSelectRegion(ctx context.Context, log zerolog.Logger, id, providerID string, organizationID string, platformc platform.PlatformServiceClient) *platform.Region {
	region, err := SelectRegion(ctx, log, id, providerID, organizationID, platformc)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list regions")
	}
	return region
}

// SelectRegion fetches the region with given ID or location or returns an error if not found.
// If no ID is specified, all regions are fetched from the selected provider
// and if the list is exactly 1 long, that region is returned.
func SelectRegion(ctx context.Context, log zerolog.Logger, id, providerID string, organizationID string, platformc platform.PlatformServiceClient) (*platform.Region, error) {
	if id == "" {
		provider, err := SelectProvider(ctx, log, providerID, organizationID, platformc)
		if err != nil {
			return nil, err
		}
		list, err := platformc.ListRegions(ctx, &platform.ListRegionsRequest{ProviderId: provider.GetId(), OrganizationId: organizationID, Options: &common.ListOptions{}})
		if err != nil {
			log.Debug().Err(err).Msg("Failed to list regions")
			return nil, err
		}
		if len(list.Items) != 1 {
			log.Debug().Err(err).Msgf("You have access to %d regions. Please specify one explicitly.", len(list.Items))
			return nil, fmt.Errorf("You have access to %d regions. Please specify one explicitly.", len(list.Items))
		}
		return list.Items[0], nil
	}
	result, err := platformc.GetRegion(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) {
			// Try to lookup region by location
			provider, err := SelectProvider(ctx, log, providerID, organizationID, platformc)
			if err != nil {
				return nil, err
			}
			list, err := platformc.ListRegions(ctx, &platform.ListRegionsRequest{ProviderId: provider.GetId(), OrganizationId: organizationID, Options: &common.ListOptions{}})
			if err == nil {
				for _, x := range list.Items {
					if x.GetLocation() == id {
						return x, nil
					}
				}
			}
		}
		log.Debug().Err(err).Str("region", id).Msg("Failed to get region")
		return nil, err
	}
	return result, nil
}
