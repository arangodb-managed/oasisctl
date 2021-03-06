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
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

// MustSelectOrganization fetches the organization with given ID, name, or URL and fails if no organization is found.
// If no ID is specified, all organizations are fetched and if the user
// is member of exactly 1, that organization is returned.
func MustSelectOrganization(ctx context.Context, log zerolog.Logger, id string, rmc rm.ResourceManagerServiceClient) *rm.Organization {
	org, err := SelectOrganization(ctx, log, id, rmc)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list organizations")
	}
	return org
}

// SelectOrganization fetches the organization with given ID, name, or URL or returns an error if not found.
// If no ID is specified, all organizations are fetched and if the user
// is member of exactly 1, that organization is returned.
func SelectOrganization(ctx context.Context, log zerolog.Logger, id string, rmc rm.ResourceManagerServiceClient) (*rm.Organization, error) {
	if id == "" {
		list, err := rmc.ListOrganizations(ctx, &common.ListOptions{})
		if err != nil {
			log.Debug().Err(err).Msg("Failed to list organizations")
			return nil, err
		}
		if len(list.Items) != 1 {
			log.Debug().Err(err).Msgf("You're member of %d organizations. Please specify one explicitly.", len(list.Items))
			return nil, fmt.Errorf("You're member of %d organizations. Please specify one explicitly.", len(list.Items))
		}
		return list.Items[0], nil
	}
	result, err := rmc.GetOrganization(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) {
			// Try to lookup organization by name or URL
			list, err := rmc.ListOrganizations(ctx, &common.ListOptions{})
			if err == nil {
				for _, x := range list.Items {
					if x.GetName() == id || x.GetUrl() == id {
						return x, nil
					}
				}
			}
		}
		log.Debug().Err(err).Str("organization", id).Msg("Failed to get organization")
		return nil, err
	}
	return result, nil
}
