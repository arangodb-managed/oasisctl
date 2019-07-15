//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Robert Stam
//

package selection

import (
	"context"

	"github.com/rs/zerolog"

	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

// MustSelectMember fetches the member with given ID.
// If no ID is specified, all members are fetched from the selected organization
// and if the list is exactly 1 long, that group is returned.
func MustSelectMember(ctx context.Context, log zerolog.Logger, id, orgID string, iamc iam.IAMServiceClient, rmc rm.ResourceManagerServiceClient) *iam.User {
	if id == "" {
		org := MustSelectOrganization(ctx, log, orgID, rmc)
		list, err := rmc.ListOrganizationMembers(ctx, &common.ListOptions{ContextId: org.GetId()})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to list organization members")
		}
		if len(list.Items) != 1 {
			log.Fatal().Err(err).Msgf("Organization contains %d members. Please specify one explicitly.", len(list.Items))
		}
		id = list.Items[0].UserId
	}
	result, err := iamc.GetUser(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) || common.IsPermissionDenied(err) {
			// Try to lookup group by name or URL
			org := MustSelectOrganization(ctx, log, orgID, rmc)
			list, err := rmc.ListOrganizationMembers(ctx, &common.ListOptions{ContextId: org.GetId()})
			if err == nil {
				for _, x := range list.Items {
					u, err := iamc.GetUser(ctx, &common.IDOptions{Id: x.GetUserId()})
					if err == nil {
						if u.GetName() == id || u.Email == id {
							return u
						}
					}
				}
			}
		}
		log.Fatal().Err(err).Str("user", id).Msg("Failed to get user")
	}
	return result
}
