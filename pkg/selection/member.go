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
// Author Robert Stam
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

// MustSelectMember fetches the member with given ID, name, or email and fails if no member is found.
// If no ID is specified, all members are fetched from the selected organization
// and if the list is exactly 1 long, that user is returned.
func MustSelectMember(ctx context.Context, log zerolog.Logger, id, orgID string, iamc iam.IAMServiceClient, rmc rm.ResourceManagerServiceClient) *iam.User {
	member, err := SelectMember(ctx, log, id, orgID, iamc, rmc)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get member")
	}
	return member
}

// SelectMember fetches the member with given ID, name, or email or returns an error if not found.
// If no ID is specified, all members are fetched from the selected organization
// and if the list is exactly 1 long, that user is returned.
func SelectMember(ctx context.Context, log zerolog.Logger, id, orgID string, iamc iam.IAMServiceClient, rmc rm.ResourceManagerServiceClient) (*iam.User, error) {
	if id == "" {
		org, err := SelectOrganization(ctx, log, orgID, rmc)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to list organization members")
			return nil, err
		}
		list, err := rmc.ListOrganizationMembers(ctx, &common.ListOptions{ContextId: org.GetId()})
		if err != nil {
			return nil, err
		}
		if len(list.Items) != 1 {
			log.Debug().Err(err).Msgf("Organization contains %d members. Please specify one explicitly.", len(list.Items))
			return nil, fmt.Errorf("Organization contains %d members. Please specify one explicitly.", len(list.Items))
		}
		id = list.Items[0].UserId
	}
	result, err := iamc.GetUser(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) || common.IsPermissionDenied(err) {
			// Try to lookup group by name or email
			org, err := SelectOrganization(ctx, log, orgID, rmc)
			if err != nil {
				return nil, err
			}
			list, err := rmc.ListOrganizationMembers(ctx, &common.ListOptions{ContextId: org.GetId()})
			if err == nil {
				for _, x := range list.Items {
					u, err := iamc.GetUser(ctx, &common.IDOptions{Id: x.GetUserId()})
					if err == nil {
						if u.GetName() == id || u.Email == id {
							return u, nil
						}
					}
				}
			}
		}
		log.Debug().Err(err).Str("user", id).Msg("Failed to get user")
		return nil, err
	}
	return result, nil
}
