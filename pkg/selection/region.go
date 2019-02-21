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
	platform "github.com/arangodb-managed/apis/platform/v1"
)

// MustSelectRegion fetches the region with given ID.
// If no ID is specified, all regions are fetched from the selected provider
// and if the list is exactly 1 long, that region is returned.
func MustSelectRegion(ctx context.Context, log zerolog.Logger, id, providerID string, platformc platform.PlatformServiceClient) *platform.Region {
	if id == "" {
		provider := MustSelectProvider(ctx, log, providerID, platformc)
		list, err := platformc.ListRegions(ctx, &common.ListOptions{ContextId: provider.GetId()})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to list regions")
		}
		if len(list.Items) != 1 {
			log.Fatal().Err(err).Msgf("You have access to %d regions. Please specify one explicitly.", len(list.Items))
		}
		return list.Items[0]
	}
	result, err := platformc.GetRegion(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) || common.IsPermissionDenied(err) {
			// Try to lookup region by name or URL
			provider := MustSelectProvider(ctx, log, providerID, platformc)
			list, err := platformc.ListRegions(ctx, &common.ListOptions{ContextId: provider.GetId()})
			if err == nil {
				for _, x := range list.Items {
					if x.GetLocation() == id {
						return x
					}
				}
			}
		}
		log.Fatal().Err(err).Str("region", id).Msg("Failed to get region")
	}
	return result
}
