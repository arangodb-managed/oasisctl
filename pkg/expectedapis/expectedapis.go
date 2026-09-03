//
// DISCLAIMER
//
// Copyright 2026 ArangoDB GmbH, Cologne, Germany
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

package expectedapis

import (
	"fmt"

	backup "github.com/arangodb-managed/apis/backup/v1"
	common "github.com/arangodb-managed/apis/common/v1"
	crypto "github.com/arangodb-managed/apis/crypto/v1"
	data "github.com/arangodb-managed/apis/data/v1"
	example "github.com/arangodb-managed/apis/example/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"
	mon "github.com/arangodb-managed/apis/monitoring/v1"
	platform "github.com/arangodb-managed/apis/platform/v1"
	replication "github.com/arangodb-managed/apis/replication/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
	security "github.com/arangodb-managed/apis/security/v1"
	tools "github.com/arangodb-managed/apis/tools/v1"
)

type apiVersion struct {
	id                string
	major, minor, patch int
}

func expectedAPIRegistry() []apiVersion {
	return []apiVersion{
		{backup.APIID, backup.APIMajorVersion, backup.APIMinorVersion, backup.APIPatchVersion},
		{crypto.APIID, crypto.APIMajorVersion, crypto.APIMinorVersion, crypto.APIPatchVersion},
		{data.APIID, data.APIMajorVersion, data.APIMinorVersion, data.APIPatchVersion},
		{example.APIID, example.APIMajorVersion, example.APIMinorVersion, example.APIPatchVersion},
		{iam.APIID, iam.APIMajorVersion, iam.APIMinorVersion, iam.APIPatchVersion},
		{mon.APIID, mon.APIMajorVersion, mon.APIMinorVersion, mon.APIPatchVersion},
		{platform.APIID, platform.APIMajorVersion, platform.APIMinorVersion, platform.APIPatchVersion},
		{replication.APIID, replication.APIMajorVersion, replication.APIMinorVersion, replication.APIPatchVersion},
		{rm.APIID, rm.APIMajorVersion, rm.APIMinorVersion, rm.APIPatchVersion},
		{security.APIID, security.APIMajorVersion, security.APIMinorVersion, security.APIPatchVersion},
	}
}

func versionString(major, minor, patch int) string {
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}

// ExpectedAPIVersions returns apiID -> semver for every API oasisctl expects.
func ExpectedAPIVersions() map[string]string {
	registry := expectedAPIRegistry()
	versions := make(map[string]string, len(registry))
	for _, api := range registry {
		versions[api.id] = versionString(api.major, api.minor, api.patch)
	}
	return versions
}

// ExpectedAPIVersionPairs returns gRPC pairs for the tools compatibility check.
func ExpectedAPIVersionPairs() []*tools.APIVersionPair {
	registry := expectedAPIRegistry()
	pairs := make([]*tools.APIVersionPair, 0, len(registry))
	for _, api := range registry {
		pairs = append(pairs, &tools.APIVersionPair{
			ApiId: api.id,
			Version: &common.Version{
				Major: int32(api.major),
				Minor: int32(api.minor),
				Patch: int32(api.patch),
			},
		})
	}
	return pairs
}
