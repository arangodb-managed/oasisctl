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
	backup "github.com/arangodb-managed/apis/backup/v1"
)

// BackupList returns a list of backups for a deployment.
func BackupList(list []*backup.Backup, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			{"id", x.Id},
			{"backup-policy-id", x.BackupPolicyId},
			{"deleted", x.IsDeleted},
			{"deployment-id", x.DeploymentId},
			{"description", x.Description},
			{"name", x.Name},
			{"upload", x.Upload},
			{"url", x.Url},
			{"state", x.Status.State},
			{"db-servers", x.DeploymentInfo.Servers.Dbservers},
			{"uploaded", x.Status.GetUploadStatus().GetUploaded()},
			{"auto-deleted-at", formatTime(opts, x.AutoDeletedAt)},
			{"created-at", formatTime(opts, x.CreatedAt)},
			{"deleted-at", formatTime(opts, x.DeletedAt)},
		}
	}, false)
}
