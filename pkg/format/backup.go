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

// Backup returns a single backup formatted for humans.
func Backup(x *backup.Backup, opts Options) string {

	data := []kv{
		{"id", x.Id},
		{"backup-policy-id", x.BackupPolicyId},
		{"deleted", x.IsDeleted},
		{"deployment-id", x.DeploymentId},
		{"description", x.Description},
		{"name", x.Name},
		{"upload", x.Upload},
		{"url", x.Url},
		{"auto-deleted-at", formatTime(opts, x.AutoDeletedAt)},
		{"created-at", formatTime(opts, x.CreatedAt)},
		{"deleted-at", formatTime(opts, x.DeletedAt)},
	}

	if x.Status != nil {
		data = append(data, kv{"state", x.Status.State}, kv{"uploaded", x.Status.GetUploadStatus().GetUploaded()})
	}
	if x.DeploymentInfo.Servers != nil {
		data = append(data, kv{"dbservers", x.DeploymentInfo.Servers.Dbservers})
	}
	return formatObject(opts, data...)
}
