//
// DISCLAIMER
//
// Copyright 2022 ArangoDB GmbH, Cologne, Germany
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

package data

import (
	"fmt"

	backup "github.com/arangodb-managed/apis/backup/v1"
	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

var copyBackupCmd = cmd.InitCommand(
	cmd.BackupCmd,
	&cobra.Command{
		Use:   "copy",
		Short: "Copy a backup from source backup to given region",
	},
	func(c *cobra.Command, f *flag.FlagSet) {
		cargs := &struct {
			sourceBackupID string
			regionID       string
		}{}
		f.StringVar(&cargs.sourceBackupID, "source-backup-id", "", "Identifier of the source backup")
		f.StringVar(&cargs.regionID, "region-id", "", "Identifier of the region where the new backup is to be created")

		c.Run = func(c *cobra.Command, args []string) {
			// Validate arguments
			log := cmd.CLILog
			sourceBackupID, argsUsed := cmd.ReqOption("source-backup-id", cargs.sourceBackupID, args, 0)
			regionID, argsUsed := cmd.ReqOption("region-id", cargs.regionID, args, 0)
			cmd.MustCheckNumberOfArgs(args, argsUsed)

			// Connect
			conn := cmd.MustDialAPI()
			backupc := backup.NewBackupServiceClient(conn)
			ctx := cmd.ContextWithToken()

			// Copy backup
			b, err := backupc.CopyBackup(ctx, &backup.CopyBackupRequest{SourceBackupId: sourceBackupID, RegionId: regionID})
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to copy backup")
			}

			// Show result
			fmt.Println(format.Backup(b, cmd.RootArgs.Format))
		}
	},
)
