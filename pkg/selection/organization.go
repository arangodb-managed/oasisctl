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
	"errors"
	"fmt"
	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
	"github.com/rs/zerolog"
)

// MustSelectOrganization fetches the organization with given ID and fails if no organization is found.
// If no ID is specified, all organizations are fetched and if the user
// is member of exactly 1, that organization is returned.
func MustSelectOrganization(ctx context.Context, log zerolog.Logger, id string, rmc rm.ResourceManagerServiceClient) *rm.Organization {

	err, org := SelectOrganization(ctx, log, id, rmc)
	if err != nil {
		log.Fatal().Err(err)
	}
	return org
}

// SelectOrganization fetches the organization with given ID or returns an error if not found.
// If no ID is specified, all organizations are fetched and if the user
// is member of exactly 1, that organization is returned.
func SelectOrganization(ctx context.Context, log zerolog.Logger, id string, rmc rm.ResourceManagerServiceClient) (error, *rm.Organization) {
	if id == "" {
		list, err := rmc.ListOrganizations(ctx, &common.ListOptions{})
		if err != nil {
			return errors.New("Failed to list organizations"), nil
		}
		if len(list.Items) != 1 {
			return errors.New(fmt.Sprintf("You're member of %d organizations. Please specify one explicitly.", len(list.Items))), nil
		}
		return nil, list.Items[0]
	}
	result, err := rmc.GetOrganization(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) {
			// Try to lookup organization by name or URL
			list, err := rmc.ListOrganizations(ctx, &common.ListOptions{})
			if err == nil {
				for _, x := range list.Items {
					if x.GetName() == id || x.GetUrl() == id {
						return nil, x
					}
				}
			}
		}
		return errors.New(err.Error() + "; organization: Failed to get organization"), nil
	}
	return nil, result
}
