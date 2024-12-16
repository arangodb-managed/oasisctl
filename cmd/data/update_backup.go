//
// DISCLAIMER
//
// Copyright 2020-2024 ArangoDB GmbH, Cologne, Germany
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
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	backup "github.com/arangodb-managed/apis/backup/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

var updateBackupCmd = cmd.InitCommand(
	cmd.UpdateCmd,
	&cobra.Command{
		Use:   "backup",
		Short: "Update a backup",
	},
	func(c *cobra.Command, f *flag.FlagSet) {
		cargs := &struct {
			backupID      string
			name          string
			description   string
			autoDeletedAt int
			upload        bool
		}{}
		f.StringVarP(&cargs.backupID, "backup-id", "d", "", "Identifier of the backup")
		f.StringVar(&cargs.name, "name", "", "Name of the backup")
		f.StringVar(&cargs.description, "description", "", "Description of the backup")
		f.BoolVar(&cargs.upload, "upload", false, "The backups should be uploaded")
		f.IntVar(&cargs.autoDeletedAt, "auto-deleted-at", 0, "Time (h) until auto delete of the backup")

		c.Run = func(c *cobra.Command, args []string) {
			// Validate arguments
			log := cmd.CLILog
			backupID, argsUsed := cmd.OptOption("backup-id", cargs.backupID, args, 0)
			cmd.MustCheckNumberOfArgs(args, argsUsed)

			// Connect
			conn := cmd.MustDialAPI()
			backupc := backup.NewBackupServiceClient(conn)
			ctx := cmd.ContextWithToken()

			// Select a backup to update
			item := selection.MustSelectBackup(ctx, log, backupID, backupc)

			// Set changes
			f := c.Flags()
			hasChanges := false
			if f.Changed("name") {
				item.Name = cargs.name
				hasChanges = true
			}
			if f.Changed("description") {
				item.Description = cargs.description
				hasChanges = true
			}
			if f.Changed("upload") {
				item.Upload = cargs.upload
				hasChanges = true
			}
			if !item.Upload && cargs.autoDeletedAt == 0 {
				cargs.autoDeletedAt = 6
				f.AddFlag(&flag.Flag{Name: "auto-deleted-at", Changed: true})
			}
			if f.Changed("auto-deleted-at") {
				t := time.Now().Add(time.Duration(cargs.autoDeletedAt) * time.Hour)
				tp := timestamppb.New(t)
				if err := tp.CheckValid(); err != nil {
					log.Fatal().Err(err).Msg("Failed to convert from time to proto time")
				}
				item.AutoDeletedAt = tp
				hasChanges = true
			}

			if !hasChanges {
				fmt.Println("No changes")
				return
			}

			// Update backup
			updated, err := backupc.UpdateBackup(ctx, item)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to update backup")
			}

			// Show result
			fmt.Println("Updated backup!")
			fmt.Println(format.Backup(updated, cmd.RootArgs.Format))
		}
	},
)
