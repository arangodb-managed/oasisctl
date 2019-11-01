//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
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
		{"state", x.Status.State},
		{"autodeletedat", formatTime(opts, x.AutoDeletedAt)},
		{"createdat", formatTime(opts, x.CreatedAt)},
		{"deletedat", formatTime(opts, x.DeletedAt)},
	}
	return formatObject(opts, data...)
}
