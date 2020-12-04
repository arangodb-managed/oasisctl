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
	metrics "github.com/arangodb-managed/apis/metrics/v1"
)

// MetricsToken returns a single metrics token formatted for humans.
func MetricsToken(x *metrics.Token, opts Options) string {
	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"token", x.GetToken()},
		kv{"name", x.GetName()},
		kv{"description", x.GetDescription()},
		kv{"revoked", formatBool(opts, x.GetIsRevoked())},
		kv{"created-at", formatTime(opts, x.GetCreatedAt())},
		kv{"expires-at", formatTime(opts, x.GetExpiresAt())},
	)
}

// MetricsTokenList returns a list of metrics token formatted for humans.
func MetricsTokenList(list []*metrics.Token, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			{"id", x.GetId()},
			{"name", x.GetName()},
			{"description", x.GetDescription()},
			{"revoked", formatBool(opts, x.GetIsRevoked())},
			{"created-at", formatTime(opts, x.GetCreatedAt())},
			{"expires-at", formatTime(opts, x.GetExpiresAt())},
		}
	}, false)
}
