//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package data

import (
	"fmt"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	backup "github.com/arangodb-managed/apis/backup/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/format"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

func init() {
	cmd.InitCommand(
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
			}{}
			f.StringVarP(&cargs.backupID, "backup-id", "d", "", "Identifier of the backup")
			f.StringVar(&cargs.name, "name", "", "Name of the backup")
			f.StringVar(&cargs.description, "description", "", "Description of the backup")
			f.IntVar(&cargs.autoDeletedAt, "autodeletedat", 6, "Time (h) until auto delete of the backup")

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
				if f.Changed("autodeletedat") {
					t := time.Now().Add(time.Duration(cargs.autoDeletedAt) * time.Hour)
					tp, err := types.TimestampProto(t)
					if err != nil {
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
}
