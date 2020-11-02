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
	audit "github.com/arangodb-managed/apis/audit/v1"
)

// AuditLogArchive returns a formatted auditlog archive.
func AuditLogArchive(x *audit.AuditLogArchive, opts Options) string {
	return formatObject(opts, generateKeyValuePairs(x, opts)...)
}

// AuditLogArchiveList returns a formatted list of auditlog archives.
func AuditLogArchiveList(list []*audit.AuditLogArchive, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return generateKeyValuePairs(x, opts)
	}, false)
}

func generateKeyValuePairs(x *audit.AuditLogArchive, opts Options) []kv {
	return []kv{
		{"id", x.GetId()},
		{"url", x.GetUrl()},
		{"auditlog-id", x.GetAuditlogId()},
		{"deployment-id", x.GetDeploymentId()},
		{"created-at", formatTime(opts, x.GetCreatedAt())},
		{"deleted-at", formatTime(opts, x.GetDeletedAt(), "-")},
		{"size-in-bytes-changed-at", formatTime(opts, x.GetSizeInBytesChangedAt(), "-")},
	}
}
