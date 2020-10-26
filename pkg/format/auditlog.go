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

package format

import (
	"fmt"
	"strings"

	audit "github.com/arangodb-managed/apis/audit/v1"
)

const (
	cloud     = "cloud"
	httpsPost = "https-post"
)

// AuditLog returns a single audit log formatted for humans.
func AuditLog(x *audit.AuditLog, opts Options) string {
	d := []kv{
		{"id", x.GetId()},
		{"name", x.GetName()},
		{"url", x.GetUrl()},
		{"description", x.GetDescription()},
		{"default", formatBool(opts, x.GetIsDefault())},
		{"created-at", formatTime(opts, x.GetCreatedAt())},
		{"deleted-at", formatTime(opts, x.GetDeletedAt(), "-")},
		{"destinations", AuditLogDestinationList(x.GetDestinations(), opts)},
	}
	return formatObject(opts, d...)
}

// AuditLog returns a single audit log formatted for humans.
func AuditLogList(list []*audit.AuditLog, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			{"id", x.GetId()},
			{"name", x.GetName()},
			{"url", x.GetUrl()},
			{"description", x.GetDescription()},
			{"default", formatBool(opts, x.GetIsDefault())},
			{"created-at", formatTime(opts, x.GetCreatedAt())},
			{"deleted-at", formatTime(opts, x.GetDeletedAt(), "-")},
			{"destinations", AuditLogDestinationList(x.GetDestinations(), opts)},
		}
	}, false)
}

// AuditLogDestinationList returns a list of configured destinations.
func AuditLogDestinationList(list []*audit.AuditLog_Destination, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		formattedList := []kv{
			{"index", i},
			{"type", x.GetType()},
		}
		if x.GetType() == httpsPost {
			formattedList = append(formattedList, []kv{
				{"url", x.GetHttpPost().GetUrl()},
				{"headers", formatHeaders(opts, x.GetHttpPost().GetHeaders())},
			}...)
		}
		return formattedList
	}, true)
}

// formatHeaders returns a list of formatted headers for a destination
func formatHeaders(opts Options, list []*audit.AuditLog_Header) string {
	var headers []string
	for _, x := range list {
		headers = append(headers, fmt.Sprintf("%s=%s", x.GetKey(), x.GetValue()))
	}
	return strings.Join(headers, ",")
}
