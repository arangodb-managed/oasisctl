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

package crypto

import (
	crypto "github.com/arangodb-managed/apis/crypto/v1"
	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/tests"
)

// getCryptoClientAndProject creates a crypto client and a project for the tests to work with.
func getCryptoClientAndProject() (crypto.CryptoServiceClient, string, error) {
	org, err := tests.GetDefaultOrganization()
	if err != nil {
		return nil, "", err
	}
	proj, err := tests.GetDefaultProject(org)
	if err != nil {
		return nil, "", err
	}

	conn := cmd.MustDialAPI()
	cryptoc := crypto.NewCryptoServiceClient(conn)
	return cryptoc, proj, nil
}
