//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
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
		if cert, found := selectDefaultCertificate(list); found {
			return cert, nil
		}
		if len(list.Items) != 1 {
			log.Debug().Err(err).Msgf("You have access to %d CA certificates and no defaults were found. Please specify one explicitly.", len(list.Items))
			return nil, fmt.Errorf("You have access to %d CA certificates and no defaults were found. Please specify one explicitly.", len(list.Items))
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

// selectDefaultCertificate looks for the default certificate in a list of certificates.
func selectDefaultCertificate(list *crypto.CACertificateList) (*crypto.CACertificate, bool) {
	for _, c := range list.GetItems() {
		if c.GetIsDefault() {
			return c, true
		}
	}
	return nil, false
}
