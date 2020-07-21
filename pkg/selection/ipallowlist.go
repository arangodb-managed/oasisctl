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
	security "github.com/arangodb-managed/apis/security/v1"
)

// MustSelectIPAllowlist fetches the IP allowlist with given ID, name, or URL and fails if no organization is found.
// If no ID is specified, all IP allowlists are fetched from the selected project
// and if the list is exactly 1 long, that IP allowlist is returned.
func MustSelectIPAllowlist(ctx context.Context, log zerolog.Logger, id, projectID, orgID string, securityc security.SecurityServiceClient, rmc rm.ResourceManagerServiceClient) *security.IPAllowlist {
	ipallowlist, err := SelectIPAllowlist(ctx, log, id, projectID, orgID, securityc, rmc)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list IP allowlists")
	}
	return ipallowlist
}

// SelectIPAllowlist fetches the IP allowlist with given ID, name, or URL and fails if no allowlist is found.
// If no ID is specified, all IP allowlists are fetched from the selected project
// and if the list is exactly 1 long, that IP allowlist is returned.
func SelectIPAllowlist(ctx context.Context, log zerolog.Logger, id, projectID, orgID string, securityc security.SecurityServiceClient, rmc rm.ResourceManagerServiceClient) (*security.IPAllowlist, error) {
	if id == "" {
		project, err := SelectProject(ctx, log, projectID, orgID, rmc)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to list IP allowlists")
			return nil, err
		}
		list, err := securityc.ListIPAllowlists(ctx, &common.ListOptions{ContextId: project.GetId()})
		if err != nil {
			return nil, err
		}
		if len(list.Items) != 1 {
			log.Debug().Err(err).Msgf("You have access to %d IP allowlists. Please specify one explicitly.", len(list.Items))
			return nil, fmt.Errorf("You have access to %d IP allowlists. Please specify one explicitly.", len(list.Items))
		}
		return list.Items[0], nil
	}
	result, err := securityc.GetIPAllowlist(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) {
			// Try to lookup IP allowlist by name or URL
			project, err := SelectProject(ctx, log, projectID, orgID, rmc)
			if err != nil {
				return nil, err
			}
			list, err := securityc.ListIPAllowlists(ctx, &common.ListOptions{ContextId: project.GetId()})
			if err == nil {
				for _, x := range list.Items {
					if x.GetName() == id || x.GetUrl() == id {
						return x, nil
					}
				}
			}
		}
		log.Debug().Err(err).Str("ipallowlist", id).Msg("Failed to get IP allowlist")
		return nil, err
	}
	return result, nil
}
