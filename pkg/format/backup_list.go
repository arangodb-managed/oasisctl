//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
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
