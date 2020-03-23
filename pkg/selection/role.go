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
	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

// MustSelectRole fetches the role with given ID, name, or URL and fails if no role is found.
// If no ID is specified, all roles are fetched from the selected organization
// and if the list is exactly 1 long, that role is returned.
func MustSelectRole(ctx context.Context, log zerolog.Logger, id, orgID string, iamc iam.IAMServiceClient, rmc rm.ResourceManagerServiceClient) *iam.Role {
	role, err := SelectRole(ctx, log, id, orgID, iamc, rmc)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list roles")
	}
	return role
}

// SelectRole fetches the role with given ID, name, or URL or returns an error if not found.
// If no ID is specified, all roles are fetched from the selected organization
// and if the list is exactly 1 long, that role is returned.
func SelectRole(ctx context.Context, log zerolog.Logger, id, orgID string, iamc iam.IAMServiceClient, rmc rm.ResourceManagerServiceClient) (*iam.Role, error) {
	if id == "" {
		org, err := SelectOrganization(ctx, log, orgID, rmc)
		if err != nil {
			return nil, err
		}
		list, err := iamc.ListRoles(ctx, &common.ListOptions{ContextId: org.GetId()})
		if err != nil {
			log.Debug().Err(err).Msg("Failed to list roles")
			return nil, err
		}
		if len(list.Items) != 1 {
			log.Debug().Err(err).Msgf("You have access to %d roles. Please specify one explicitly.", len(list.Items))
			return nil, fmt.Errorf("You have access to %d roles. Please specify one explicitly.", len(list.Items))
		}
		return list.Items[0], nil
	}
	result, err := iamc.GetRole(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) {
			// Try to lookup role by name or URL
			org, err := SelectOrganization(ctx, log, orgID, rmc)
			if err != nil {
				return nil, err
			}
			list, err := iamc.ListRoles(ctx, &common.ListOptions{ContextId: org.GetId()})
			if err == nil {
				for _, x := range list.Items {
					if x.GetName() == id || x.GetUrl() == id {
						return x, nil
					}
				}
			}
		}
		log.Debug().Err(err).Str("role", id).Msg("Failed to get role")
		return nil, err
	}
	return result, nil
}
