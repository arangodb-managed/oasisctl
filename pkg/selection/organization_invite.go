//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package selection

import (
	"context"

	"github.com/rs/zerolog"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

// MustSelectOrganizationInvite fetches the organization invite with given ID.
// If no ID is specified, all invites are fetched from the selected organization
// and if the list is exactly 1 long, that invite is returned.
func MustSelectOrganizationInvite(ctx context.Context, log zerolog.Logger, id, orgID string, rmc rm.ResourceManagerServiceClient) *rm.OrganizationInvite {
	if id == "" {
		org := MustSelectOrganization(ctx, log, orgID, rmc)
		list, err := rmc.ListOrganizationInvites(ctx, &common.ListOptions{ContextId: org.GetId()})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to list organization invites")
		}
		if len(list.Items) != 1 {
			log.Fatal().Err(err).Msgf("You have access to %d organization invites. Please specify one explicitly.", len(list.Items))
		}
		return list.Items[0]
	}
	result, err := rmc.GetOrganizationInvite(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) {
			// Try to lookup organization invite by name or URL
			org := MustSelectOrganization(ctx, log, orgID, rmc)
			list, err := rmc.ListOrganizationInvites(ctx, &common.ListOptions{ContextId: org.GetId()})
			if err == nil {
				for _, x := range list.Items {
					if x.GetEmail() == id || x.GetUrl() == id {
						return x
					}
				}
			}
		}
		log.Fatal().Err(err).Str("organization-invite", id).Msg("Failed to get organization invite")
	}
	return result
}
