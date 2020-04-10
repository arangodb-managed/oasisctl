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
// Author Brautigam Gergely
//

package format

import (
	example "github.com/arangodb-managed/apis/example/v1"
)

// ExampleDatasetInstallation returns a single installation formatted for humans.
func ExampleDatasetInstallation(x *example.ExampleDatasetInstallation, opts Options) string {

	data := []kv{
		{"id", x.Id},
		{"deleted", x.IsDeleted},
		{"example-dataset-id", x.ExampledatasetId},
		{"deployment-id", x.DeploymentId},
		{"url", x.Url},
		{"created-at", formatTime(opts, x.CreatedAt)},
		{"deleted-at", formatTime(opts, x.DeletedAt)},
	}

	if x.Status != nil {
		data = append(data,
			kv{"database", x.Status.GetDatabaseName()},
			kv{"state", x.Status.GetState()},
			kv{"failed", x.Status.GetIsFailed()},
			kv{"available", x.Status.GetIsAvailable()})
	}
	return formatObject(opts, data...)
}

// ExampleDatasetInstallationList returns a list of installations formatted for humans.
func ExampleDatasetInstallationList(list []*example.ExampleDatasetInstallation, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		data := []kv{
			{"id", x.Id},
			{"deleted", x.IsDeleted},
			{"example-dataset-id", x.ExampledatasetId},
			{"deployment-id", x.DeploymentId},
			{"url", x.Url},
			{"created-at", formatTime(opts, x.CreatedAt)},
			{"deleted-at", formatTime(opts, x.DeletedAt)},
		}
		if x.Status != nil {
			data = append(data,
				kv{"database", x.Status.GetDatabaseName()},
				kv{"state", x.Status.GetState()},
				kv{"failed", x.Status.GetIsFailed()},
				kv{"available", x.Status.GetIsAvailable()})
		}
		return data
	}, false)
}
