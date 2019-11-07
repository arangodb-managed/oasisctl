//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Gergely Brautigam
//

package selection

import (
	"context"

	backup "github.com/arangodb-managed/apis/backup/v1"
	common "github.com/arangodb-managed/apis/common/v1"
	"github.com/rs/zerolog"
)

// MustSelectBackup TODO: Implement this.
func MustSelectBackup(ctx context.Context, log zerolog.Logger, id string, backupc backup.BackupServiceClient) *backup.Backup {
	if id == "" {
		list, err := backupc.ListBackups(ctx, &backup.ListBackupsRequest{DeploymentId: id})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to list backups")
		}
		if len(list.Items) != 1 {
			log.Fatal().Err(err).Msgf("You have access to %d backups. Please specify one explicitly.", len(list.Items))
		}
		return list.Items[0]
	}
	result, err := backupc.GetBackup(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) {
			// Try to lookup deployment by name or URL
			list, err := backupc.ListBackups(ctx, &backup.ListBackupsRequest{DeploymentId: id})
			if err == nil {
				for _, x := range list.Items {
					if x.GetName() == id || x.GetUrl() == id {
						return x
					}
				}
			}
		}
		log.Fatal().Err(err).Str("backup", id).Msg("Failed to get backup")
	}
	return result
}
