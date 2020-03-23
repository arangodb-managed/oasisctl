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

package format

import (
	data "github.com/arangodb-managed/apis/data/v1"
)

// Version returns a single version formatted for humans.
func Version(x *data.Version, opts Options) string {
	return formatObject(opts,
		kv{"version", x.GetVersion()},
	)
}

// VersionList returns a list of versions formatted for humans.
func VersionList(list []*data.Version, defaultVersion *data.Version, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			kv{"version", x.GetVersion()},
			kv{"default", formatBool(opts, x.GetVersion() == defaultVersion.GetVersion())},
		}
	}, true)
}
