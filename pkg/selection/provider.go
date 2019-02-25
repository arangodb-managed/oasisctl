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

// MustSelectProvider fetches the provider with given ID.
// If no ID is specified, all providers are fetched and if the user
// is member of exactly 1, that provider is returned.
func MustSelectProvider(ctx context.Context, log zerolog.Logger, id string, platformc platform.PlatformServiceClient) *platform.Provider {
	if id == "" {
		list, err := platformc.ListProviders(ctx, &common.ListOptions{})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to list providers")
		}
		if len(list.Items) != 1 {
			log.Fatal().Err(err).Msgf("You're member of %d providers. Please specify one explicitly.", len(list.Items))
		}
		return list.Items[0]
	}
	result, err := platformc.GetProvider(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) {
			// Try to lookup provider by name or URL
			list, err := platformc.ListProviders(ctx, &common.ListOptions{})
			if err == nil {
				for _, x := range list.Items {
					if x.GetName() == id {
						return x
					}
				}
			}
		}
		log.Fatal().Err(err).Str("provider", id).Msg("Failed to get provider")
	}
	return result
}
