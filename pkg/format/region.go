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
	platform "github.com/arangodb-managed/apis/platform/v1"
)

// Region returns a single region formatted for humans.
func Region(x *platform.Region, opts Options) string {
	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"provider-id", x.GetProviderId()},
		kv{"location", x.GetLocation()},
		kv{"available", formatBool(opts, x.GetAvailable())},
	)
}

// RegionList returns a list of regions formatted for humans.
func RegionList(list []*platform.Region, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			kv{"id", x.GetId()},
			kv{"provider-id", x.GetProviderId()},
			kv{"location", x.GetLocation()},
			kv{"available", formatBool(opts, x.GetAvailable())},
		}
	}, false)
}
