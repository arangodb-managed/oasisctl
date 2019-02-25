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
	data "github.com/arangodb-managed/apis/data/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

// MustSelectDeployment fetches the deployment with given ID.
// If no ID is specified, all deployments are fetched from the selected project
// and if the list is exactly 1 long, that deployment is returned.
func MustSelectDeployment(ctx context.Context, log zerolog.Logger, id, projectID, orgID string, datac data.DataServiceClient, rmc rm.ResourceManagerServiceClient) *data.Deployment {
	if id == "" {
		project := MustSelectProject(ctx, log, projectID, orgID, rmc)
		list, err := datac.ListDeployments(ctx, &common.ListOptions{ContextId: project.GetId()})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to list deployments")
		}
		if len(list.Items) != 1 {
			log.Fatal().Err(err).Msgf("You have access to %d deployments. Please specify one explicitly.", len(list.Items))
		}
		return list.Items[0]
	}
	result, err := datac.GetDeployment(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) {
			// Try to lookup deployment by name or URL
			project := MustSelectProject(ctx, log, projectID, orgID, rmc)
			list, err := datac.ListDeployments(ctx, &common.ListOptions{ContextId: project.GetId()})
			if err == nil {
				for _, x := range list.Items {
					if x.GetName() == id || x.GetUrl() == id {
						return x
					}
				}
			}
		}
		log.Fatal().Err(err).Str("deployment", id).Msg("Failed to get deployment")
	}
	return result
}
