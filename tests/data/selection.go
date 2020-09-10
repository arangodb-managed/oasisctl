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
// +build e2e

package data

import (
	"errors"

	common "github.com/arangodb-managed/apis/common/v1"
	crypto "github.com/arangodb-managed/apis/crypto/v1"
	platform "github.com/arangodb-managed/apis/platform/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
)

// getRegion returns the first available region for an organization and a provider.
func getRegion(provider, org string) (string, error) {
	ctx := cmd.ContextWithToken()
	conn := cmd.MustDialAPI()
	platformc := platform.NewPlatformServiceClient(conn)
	// select the first available region
	list, err := platformc.ListRegions(ctx, &platform.ListRegionsRequest{OrganizationId: org, ProviderId: provider})
	if err != nil {
		return "", err
	}
	if len(list.GetItems()) == 0 {
		return "", errors.New("region list is empty")
	}
	return list.GetItems()[0].GetId(), nil
}

func getDefaultCertificate(proj string) (*crypto.CACertificate, error) {
	ctx := cmd.ContextWithToken()
	conn := cmd.MustDialAPI()
	cryptoc := crypto.NewCryptoServiceClient(conn)
	list, err := cryptoc.ListCACertificates(ctx, &common.ListOptions{ContextId: proj})
	if err != nil {
		return nil, err
	}
	if len(list.GetItems()) == 0 {
		return nil, errors.New("no certificates found for project")
	}
	for _, c := range list.GetItems() {
		if c.GetIsDefault() {
			return c, nil
		}
	}
	return nil, errors.New("no default certificate found")
}

func getTermsAndConditions(org string) (string, error) {
	ctx := cmd.ContextWithToken()
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	tandc, err := rmc.GetCurrentTermsAndConditions(ctx, &common.IDOptions{Id: org})
	if err != nil {
		return "", err
	}
	return tandc.GetId(), nil
}
