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

// AuditLogEventList returns a formatted list of auditlog events.
func AuditLogEventList(list []*audit.AuditLogEvent, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			{"timestamp", formatTime(opts, x.GetTimestamp())},
			{"auditlog-archive-id", x.GetAuditlogarchiveId()},
			{"project-id", x.GetProjectId()},
			{"instance-id", x.GetInstanceId()},
			{"server", x.GetServerId()},
			{"topic", x.GetTopic()},
			{"username", x.GetUserId()},
			{"database", x.GetDatabase()},
			{"client-ip", x.GetClientIp()},
			{"authentication", x.GetAuthentication()},
			{"message", x.GetMessage()},
		}
	}, false)
}
