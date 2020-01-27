//
// DISCLAIMER
//
// Copyright 2019 ArangoDB GmbH, Cologne, Germany
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

// MustSelectOrganizationInvite fetches the organization invite with given ID, email, or URL and fails if no invite is found.
// If no ID is specified, all invites are fetched from the selected organization
// and if the list is exactly 1 long, that invite is returned.
func MustSelectOrganizationInvite(ctx context.Context, log zerolog.Logger, id, orgID string, rmc rm.ResourceManagerServiceClient) *rm.OrganizationInvite {
	invite, err := SelectOrganizationInvite(ctx, log, id, orgID, rmc)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get organization invites")
	}
	return invite
}

// SelectOrganizationInvite fetches the organization invite with given ID, email, or URL or returns an error if not found.
// If no ID is specified, all invites are fetched from the selected organization
// and if the list is exactly 1 long, that invite is returned.
func SelectOrganizationInvite(ctx context.Context, log zerolog.Logger, id, orgID string, rmc rm.ResourceManagerServiceClient) (*rm.OrganizationInvite, error) {
	if id == "" {
		org, err := SelectOrganization(ctx, log, orgID, rmc)
		if err != nil {
			return nil, err
		}
		list, err := rmc.ListOrganizationInvites(ctx, &common.ListOptions{ContextId: org.GetId()})
		if err != nil {
			log.Debug().Err(err).Msg("Failed to list organization invites")
			return nil, err
		}
		if len(list.Items) != 1 {
			log.Debug().Err(err).Msgf("You have access to %d organization invites. Please specify one explicitly.", len(list.Items))
			return nil, fmt.Errorf("You have access to %d organization invites. Please specify one explicitly.", len(list.Items))
		}
		return list.Items[0], nil
	}
	result, err := rmc.GetOrganizationInvite(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) {
			// Try to lookup organization invite by email or URL
			org, err := SelectOrganization(ctx, log, orgID, rmc)
			if err != nil {
				return nil, err
			}
			list, err := rmc.ListOrganizationInvites(ctx, &common.ListOptions{ContextId: org.GetId()})
			if err == nil {
				for _, x := range list.Items {
					if x.GetEmail() == id || x.GetUrl() == id {
						return x, nil
					}
				}
			}
		}
		log.Debug().Err(err).Str("organization-invite", id).Msg("Failed to get organization invite")
		return nil, err
	}
	return result, nil
}
