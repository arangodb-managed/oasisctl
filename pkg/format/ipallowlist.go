//
// DISCLAIMER
//
// Copyright 2020-2021 ArangoDB GmbH, Cologne, Germany
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
	"strings"

	security "github.com/arangodb-managed/apis/security/v1"
)

// IPAllowlist returns a single IP allowlist formatted for humans.
func IPAllowlist(x *security.IPAllowlist, opts Options) string {
	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"name", x.GetName()},
		kv{"description", x.GetDescription()},
		kv{"cidr-ranges", strings.Join(x.GetCidrRanges(), ", ")},
		kv{"remote-inspection-allowed", formatBool(opts, x.GetRemoteInspectionAllowed())},
		kv{"url", x.GetUrl()},
		kv{"locked", formatBool(opts, x.GetLocked())},
		kv{"created-at", formatTime(opts, x.GetCreatedAt())},
	)
}

// IPAllowlistList returns a list of IP allowlists formatted for humans.
func IPAllowlistList(list []*security.IPAllowlist, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			{"id", x.GetId()},
			{"name", x.GetName()},
			{"description", x.GetDescription()},
			{"cidr-ranges", strings.Join(x.GetCidrRanges(), ", ")},
			{"remote-inspection-allowed", formatBool(opts, x.GetRemoteInspectionAllowed())},
			{"url", x.GetUrl()},
			{"locked", formatBool(opts, x.GetLocked())},
			{"created-at", formatTime(opts, x.GetCreatedAt())},
		}
	}, false)
}
