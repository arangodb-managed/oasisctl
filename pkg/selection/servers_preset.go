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

	data "github.com/arangodb-managed/apis/data/v1"
)

// MustSelectServersSpec fetches the servers spec with given name.
func MustSelectServersSpec(ctx context.Context, log zerolog.Logger, name, projectID, regionID string, datac data.DataServiceClient) *data.Deployment_ServersSpec {
	list, err := datac.ListServersSpecPresets(ctx, &data.ServersSpecPresetsRequest{
		ProjectId: projectID,
		RegionId:  regionID,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get servers preset")
	}
	for _, sp := range list.Items {
		if sp.GetName() == name {
			return sp.GetServers()
		}
	}

	log.Fatal().Str("servers-preset", name).Msg("Failed to get servers preset: not found")
	return nil
}
