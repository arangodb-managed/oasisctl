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
	"fmt"

	"github.com/rs/zerolog"

	backup "github.com/arangodb-managed/apis/backup/v1"
	common "github.com/arangodb-managed/apis/common/v1"
)

// MustSelectBackup fetches a backup with given ID, name, or URL and fails if no backup is found.
// If no ID is specified, all backups are fetched from the selected deployment
// and if the list is exactly 1 long, that backup is returned.
func MustSelectBackup(ctx context.Context, log zerolog.Logger, id string, backupc backup.BackupServiceClient) *backup.Backup {
	backup, err := SelectBackup(ctx, log, id, backupc)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list backup")
	}
	return backup
}

// SelectBackup fetches a backup with given ID, name, or URL or returns an error if not found.
// If no ID is specified, all backups are fetched from the selected deployment
// and if the list is exactly 1 long, that backup is returned.
func SelectBackup(ctx context.Context, log zerolog.Logger, id string, backupc backup.BackupServiceClient) (*backup.Backup, error) {
	if id == "" {
		list, err := backupc.ListBackups(ctx, &backup.ListBackupsRequest{DeploymentId: id})
		if err != nil {
			log.Debug().Err(err).Msg("Failed to list backups")
			return nil, err
		}
		if len(list.Items) != 1 {
			log.Debug().Err(err).Msgf("You have access to %d backups. Please specify one explicitly.", len(list.Items))
			return nil, fmt.Errorf("You have access to %d backups. Please specify one explicitly.", len(list.Items))
		}
		return list.Items[0], nil
	}
	result, err := backupc.GetBackup(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) {
			// Try to lookup deployment by name or URL
			list, err := backupc.ListBackups(ctx, &backup.ListBackupsRequest{DeploymentId: id})
			if err == nil {
				for _, x := range list.Items {
					if x.GetName() == id || x.GetUrl() == id {
						return x, nil
					}
				}
			}
		}
		log.Debug().Err(err).Str("backup", id).Msg("Failed to get backup")
		return nil, err
	}
	return result, nil
}
