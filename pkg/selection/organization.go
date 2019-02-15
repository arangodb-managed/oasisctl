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

// MustSelectOrganization fetches the organization with given ID.
// If no ID is specified, all organizations are fetched and if the user
// is member of exactly 1, that organization is returned.
func MustSelectOrganization(ctx context.Context, log zerolog.Logger, id string, rmc rm.ResourceManagerServiceClient) *rm.Organization {
	if id == "" {
		list, err := rmc.ListOrganizations(ctx, &common.ListOptions{})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to list organizations")
		}
		if len(list.Items) != 1 {
			log.Fatal().Err(err).Msgf("You're member of %d organizations. Please specify one explicitly.", len(list.Items))
		}
		return list.Items[0]
	}
	result, err := rmc.GetOrganization(ctx, &common.IDOptions{Id: id})
	if err != nil {
		log.Fatal().Err(err).Str("organization", id).Msg("Failed to get organization")
	}
	return result
}
