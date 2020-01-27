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
	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

// MustSelectGroup fetches the group with given ID, name, or URL and fails if no group is found.
// If no ID is specified, all groups are fetched from the selected organization
// and if the list is exactly 1 long, that group is returned.
func MustSelectGroup(ctx context.Context, log zerolog.Logger, id, orgID string, iamc iam.IAMServiceClient, rmc rm.ResourceManagerServiceClient) *iam.Group {
	group, err := SelectGroup(ctx, log, id, orgID, iamc, rmc)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list groups")
	}
	return group
}

// SelectGroup fetches the group  with given ID, name, or URL or returns an error if not found.
// If no ID is specified, all groups are fetched from the selected organization
// and if the list is exactly 1 long, that group is returned.
func SelectGroup(ctx context.Context, log zerolog.Logger, id, orgID string, iamc iam.IAMServiceClient, rmc rm.ResourceManagerServiceClient) (*iam.Group, error) {
	if id == "" {
		org, err := SelectOrganization(ctx, log, orgID, rmc)
		if err != nil {
			return nil, err
		}
		list, err := iamc.ListGroups(ctx, &common.ListOptions{ContextId: org.GetId()})
		if err != nil {
			log.Debug().Err(err).Msg("Failed to list groups")
			return nil, err
		}
		if len(list.Items) != 1 {
			log.Debug().Err(err).Msgf("You have access to %d groups. Please specify one explicitly.", len(list.Items))
			return nil, fmt.Errorf(" have access to %d groups. Please specify one explicitly.", len(list.Items))
		}
		return list.Items[0], nil
	}
	result, err := iamc.GetGroup(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) {
			// Try to lookup group by name or URL
			org, err := SelectOrganization(ctx, log, orgID, rmc)
			if err != nil {
				return nil, err
			}
			list, err := iamc.ListGroups(ctx, &common.ListOptions{ContextId: org.GetId()})
			if err == nil {
				for _, x := range list.Items {
					if x.GetName() == id || x.GetUrl() == id {
						return x, nil
					}
				}
			}
		}
		log.Debug().Err(err).Str("group", id).Msg("Failed to get group")
		return nil, err
	}
	return result, nil
}
