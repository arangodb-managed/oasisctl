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

// MustSelectProject fetches the project with given ID.
// If no ID is specified, all projects are fetched from the selected organization
// and if the list is exactly 1 long, that project is returned.
func MustSelectProject(ctx context.Context, log zerolog.Logger, id, orgID string, rmc rm.ResourceManagerServiceClient) *rm.Project {
	if id == "" {
		org := MustSelectOrganization(ctx, log, orgID, rmc)
		list, err := rmc.ListProjects(ctx, &common.ListOptions{ContextId: org.GetId()})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to list projects")
		}
		if len(list.Items) != 1 {
			log.Fatal().Err(err).Msgf("You have access to %d projects. Please specify one explicitly.", len(list.Items))
		}
		return list.Items[0]
	}
	result, err := rmc.GetProject(ctx, &common.IDOptions{Id: id})
	if err != nil {
		log.Fatal().Err(err).Str("project", id).Msg("Failed to get project")
	}
	return result
}
