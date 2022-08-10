//
// DISCLAIMER
//
// Copyright 2021-2022 ArangoDB GmbH, Cologne, Germany
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

package format

import (
	"fmt"
	"strings"

	nw "github.com/arangodb-managed/apis/network/v1"
)

// PrivateEndpointService returns a single private endpoint service formatted for humans.
func PrivateEndpointService(x *nw.PrivateEndpointService, opts Options) string {
	obj := []kv{
		{"id", x.GetId()},
		{"name", x.GetName()},
		{"description", x.GetDescription()},
		{"url", x.GetUrl()},
		{"alt-dns-names", formatOptionalString(strings.Join(x.GetAlternateDnsNames(), ", "))},

		{"ready", formatBool(opts, x.GetStatus().GetReady())},
		{"needs-attention", formatBool(opts, x.GetStatus().GetNeedsAttention())},
		{"message", formatOptionalString(x.GetStatus().GetMessage())},

		{"created-at", formatTime(opts, x.GetCreatedAt())},
		{"deleted-at", formatTime(opts, x.GetDeletedAt(), "-")},
	}
	// AKS settings (if any)
	if aks := x.GetAks(); aks != nil {
		obj = append(obj,
			kv{"client-subscription-ids", formatOptionalString(strings.Join(aks.GetClientSubscriptionIds(), ", "))},
			kv{"azure-alias", formatOptionalString(x.GetStatus().GetAks().GetAlias())},
			kv{"azure-private-endpoints", len(x.GetStatus().GetAks().GetPrivateEndpointConnections())})
	}
	// AWS settings (if any)
	if aws := x.GetAws(); aws != nil {
		for _, p := range aws.GetAwsPrincipals() {
			obj = append(obj,
				kv{"aws-principals", fmt.Sprintf("AccountID=%s (Roles=%s; Users=%s)", p.GetAccountId(), strings.Join(p.GetRoleNames(), ", "), strings.Join(p.GetUserNames(), ", "))})
		}
		obj = append(obj,
			kv{"aws-service-name", formatOptionalString(x.GetStatus().GetAws().GetServiceName())},
			kv{"aws-availability-zones", formatOptionalString(strings.Join(x.GetStatus().GetAws().GetAvailabilityZones(), ", "))},
			kv{"aws-private-endpoints", len(x.GetStatus().GetAws().GetPrivateEndpointConnections())})
	}
	// GCP settings (if any)
	if gcp := x.GetGcp(); gcp != nil {
		obj = append(obj,
			kv{"gcp-projects", formatOptionalString(strings.Join(gcp.GetProjects(), ", "))},
			kv{"gcp-service-attachment", formatOptionalString(x.GetStatus().GetGcp().GetServiceAttachment())},
			kv{"gcp-private-endpoints", len(x.GetStatus().GetGcp().GetPrivateEndpointConnections())})
	}
	return formatObject(opts, obj...)
}
