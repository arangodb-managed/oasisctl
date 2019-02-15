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
	crypto "github.com/arangodb-managed/apis/crypto/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

// MustSelectCACertificate fetches the CA certificate with given ID.
// If no ID is specified, all CA certificate are fetched from the selected project
// and if the list is exactly 1 long, that CA certificate is returned.
func MustSelectCACertificate(ctx context.Context, log zerolog.Logger, id, projectID, orgID string, cryptoc crypto.CryptoServiceClient, rmc rm.ResourceManagerServiceClient) *crypto.CACertificate {
	if id == "" {
		project := MustSelectProject(ctx, log, projectID, orgID, rmc)
		list, err := cryptoc.ListCACertificates(ctx, &common.ListOptions{ContextId: project.GetId()})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to list CA certificates")
		}
		if len(list.Items) != 1 {
			log.Fatal().Err(err).Msgf("You have access to %d CA certificates. Please specify one explicitly.", len(list.Items))
		}
		return list.Items[0]
	}
	result, err := cryptoc.GetCACertificate(ctx, &common.IDOptions{Id: id})
	if err != nil {
		log.Fatal().Err(err).Str("cacertificate", id).Msg("Failed to get CA certificate")
	}
	return result
}
