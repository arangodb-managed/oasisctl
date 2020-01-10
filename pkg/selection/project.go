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
	"fmt"

	"github.com/rs/zerolog"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

// MustSelectProject fetches the project with given ID, name, or URL and fails if no project is found.
// If no ID is specified, all projects are fetched from the selected organization
// and if the list is exactly 1 long, that project is returned.
func MustSelectProject(ctx context.Context, log zerolog.Logger, id, orgID string, rmc rm.ResourceManagerServiceClient) *rm.Project {
	project, err := SelectProject(ctx, log, id, orgID, rmc)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get project")
	}
	return project
}

// MustSelectProject fetches the project with given ID, name, or URL or returns an error if not found.
// If no ID is specified, all projects are fetched from the selected organization
// and if the list is exactly 1 long, that project is returned.
func SelectProject(ctx context.Context, log zerolog.Logger, id, orgID string, rmc rm.ResourceManagerServiceClient) (*rm.Project, error) {
	if id == "" {
		org, err := SelectOrganization(ctx, log, orgID, rmc)
		if err != nil {
			return nil, err
		}
		list, err := rmc.ListProjects(ctx, &common.ListOptions{ContextId: org.GetId()})
		if err != nil {
			log.Debug().Err(err).Msg("Failed to list projects")
			return nil, err
		}
		if len(list.Items) != 1 {
			log.Debug().Err(err).Msgf("You have access to %d projects. Please specify one explicitly.", len(list.Items))
			return nil, fmt.Errorf("You have access to %d projects. Please specify one explicitly.", len(list.Items))
		}
		return list.Items[0], nil
	}
	result, err := rmc.GetProject(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) {
			// Try to lookup project by name or URL
			org, err := SelectOrganization(ctx, log, orgID, rmc)
			if err != nil {
				return nil, err
			}
			list, err := rmc.ListProjects(ctx, &common.ListOptions{ContextId: org.GetId()})
			if err == nil {
				for _, x := range list.Items {
					if x.GetName() == id || x.GetUrl() == id {
						return x, nil
					}
				}
			}
		}
		log.Debug().Err(err).Str("project", id).Msg("Failed to get project")
		return nil, err
	}
	return result, nil
}
