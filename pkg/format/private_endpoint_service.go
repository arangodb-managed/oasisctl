//
// DISCLAIMER
//
// Copyright 2021 ArangoDB GmbH, Cologne, Germany
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

	nw "github.com/arangodb-managed/apis/network/v1"
)

// PrivateEndpointService returns a single private endpoint service formatted for humans.
func PrivateEndpointService(x *nw.PrivateEndpointService, opts Options) string {
	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"name", x.GetName()},
		kv{"description", x.GetDescription()},
		kv{"url", x.GetUrl()},
		kv{"alt-dns-names", formatOptionalString(strings.Join(x.GetAlternateDnsNames(), ", "))},
		kv{"client-subscription-ids", formatOptionalString(strings.Join(x.GetAks().GetClientSubscriptionIds(), ", "))},
		kv{"ready", formatBool(opts, x.GetStatus().GetReady())},
		kv{"needs-attention", formatBool(opts, x.GetStatus().GetNeedsAttention())},
		kv{"message", formatOptionalString(x.GetStatus().GetMessage())},
		kv{"created-at", formatTime(opts, x.GetCreatedAt())},
		kv{"deleted-at", formatTime(opts, x.GetDeletedAt(), "-")},
	)
}
