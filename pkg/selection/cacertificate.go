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
	crypto "github.com/arangodb-managed/apis/crypto/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

// MustSelectCACertificate fetches the CA certificate with given ID, name, or URL and fails if no certificate is found.
// If no ID is specified, all CA certificate are fetched from the selected project
// and if the list is exactly 1 long, that CA certificate is returned.
func MustSelectCACertificate(ctx context.Context, log zerolog.Logger, id, projectID, orgID string, cryptoc crypto.CryptoServiceClient, rmc rm.ResourceManagerServiceClient) *crypto.CACertificate {
	cert, err := SelectCACertificate(ctx, log, id, projectID, orgID, cryptoc, rmc)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list CA certificates")
	}
	return cert
}

// SelectCACertificate fetches the CA certificate with given ID, name, or URL and returns an error if not found..
// If no ID is specified, all CA certificate are fetched from the selected project
// and if the list is exactly 1 long, that CA certificate is returned.
func SelectCACertificate(ctx context.Context, log zerolog.Logger, id, projectID, orgID string, cryptoc crypto.CryptoServiceClient, rmc rm.ResourceManagerServiceClient) (*crypto.CACertificate, error) {
	if id == "" {
		project, err := SelectProject(ctx, log, projectID, orgID, rmc)
		if err != nil {
			return nil, err
		}
		list, err := cryptoc.ListCACertificates(ctx, &common.ListOptions{ContextId: project.GetId()})
		if err != nil {
			log.Debug().Err(err).Msg("Failed to list CA certificates")
			return nil, err
		}
		if len(list.Items) != 1 {
			log.Debug().Err(err).Msgf("You have access to %d CA certificates. Please specify one explicitly.", len(list.Items))
			return nil, fmt.Errorf("You have access to %d CA certificates. Please specify one explicitly.", len(list.Items))
		}
		return list.Items[0], nil
	}
	result, err := cryptoc.GetCACertificate(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) {
			// Try to lookup CA certificate by name or URL
			project, err := SelectProject(ctx, log, projectID, orgID, rmc)
			if err != nil {
				return nil, err
			}
			list, err := cryptoc.ListCACertificates(ctx, &common.ListOptions{ContextId: project.GetId()})
			if err == nil {
				for _, x := range list.Items {
					if x.GetName() == id || x.GetUrl() == id {
						return x, nil
					}
				}
			}
		}
		log.Debug().Err(err).Str("cacertificate", id).Msg("Failed to get CA certificate")
		return nil, err
	}
	return result, nil
}
