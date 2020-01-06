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
			return nil, err
		}
		list, err := rmc.ListOrganizationMembers(ctx, &common.ListOptions{ContextId: org.GetId()})
		if err != nil {
			return nil, err
		}
		if len(list.Items) != 1 {
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
		return nil, err
	}
	return result, nil
}
