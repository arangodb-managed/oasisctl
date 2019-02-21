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
	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

// MustSelectGroup fetches the group with given ID.
// If no ID is specified, all groups are fetched from the selected organization
// and if the list is exactly 1 long, that group is returned.
func MustSelectGroup(ctx context.Context, log zerolog.Logger, id, orgID string, iamc iam.IAMServiceClient, rmc rm.ResourceManagerServiceClient) *iam.Group {
	if id == "" {
		org := MustSelectOrganization(ctx, log, orgID, rmc)
		list, err := iamc.ListGroups(ctx, &common.ListOptions{ContextId: org.GetId()})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to list groups")
		}
		if len(list.Items) != 1 {
			log.Fatal().Err(err).Msgf("You have access to %d groups. Please specify one explicitly.", len(list.Items))
		}
		return list.Items[0]
	}
	result, err := iamc.GetGroup(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) || common.IsPermissionDenied(err) {
			// Try to lookup group by name or URL
			org := MustSelectOrganization(ctx, log, orgID, rmc)
			list, err := iamc.ListGroups(ctx, &common.ListOptions{ContextId: org.GetId()})
			if err == nil {
				for _, x := range list.Items {
					if x.GetName() == id || x.GetUrl() == id {
						return x
					}
				}
			}
		}
		log.Fatal().Err(err).Str("group", id).Msg("Failed to get group")
	}
	return result
}
