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
		{"autodeletedat", formatTime(opts, x.AutoDeletedAt)},
		{"createdat", formatTime(opts, x.CreatedAt)},
		{"deletedat", formatTime(opts, x.DeletedAt)},
	}

	if x.Status != nil {
		data = append(data, kv{"state", x.Status.State})
	}
	if x.DeploymentInfo.Servers != nil {
		data = append(data, kv{"dbservers", x.DeploymentInfo.Servers.Dbservers})
	}
	return formatObject(opts, data...)
}
