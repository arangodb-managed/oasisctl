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
// Author Gergely Brautigam
//

package selection

import (
	"context"

	"github.com/rs/zerolog"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

// MustSelectTermsAndConditions fetches the terms and conditions with given ID.
// If no ID is specified, the current terms and conditions is selected for a given ogranization.
func MustSelectTermsAndConditions(ctx context.Context, log zerolog.Logger, id, orgID string, rmc rm.ResourceManagerServiceClient) *rm.TermsAndConditions {
	tandc, err := SelectTermsAndConditions(ctx, log, id, orgID, rmc)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get terms and conditions")
	}
	return tandc
}

// SelectTermsAndConditions fetches the terms and conditions with given ID.
// If no ID is specified, the current terms and conditions is selected for a given ogranization.
func SelectTermsAndConditions(ctx context.Context, log zerolog.Logger, id, orgID string, rmc rm.ResourceManagerServiceClient) (*rm.TermsAndConditions, error) {
	if id == "" {
		org, err := SelectOrganization(ctx, log, orgID, rmc)
		if err != nil {
			return nil, err
		}
		return rmc.GetCurrentTermsAndConditions(ctx, &common.IDOptions{Id: org.GetId()})
	}
	return rmc.GetTermsAndConditions(ctx, &common.IDOptions{Id: id})
}
