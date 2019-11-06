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
				backupID string
				name     string
			}{}
			f.StringVarP(&cargs.backupID, "backup-id", "d", "", "Identifier of the deployment")
			f.StringVar(&cargs.name, "name", "", "Name of the deployment")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				backupID, argsUsed := cmd.OptOption("deployment-id", cargs.backupID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				backupc := backup.NewBackupServiceClient(conn)
				ctx := cmd.ContextWithToken()

				item := selection.MustSelectBackup(backupID)

				// Set changes
				f := c.Flags()
				hasChanges := false
				if f.Changed("name") {
					item.Name = cargs.name
					hasChanges = true
				}
				if !hasChanges {
					fmt.Println("No changes")
					return
				}

				// Update deployment
				updated, err := backupc.UpdateBackup(ctx, item)
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to update deployment")
				}

				// Show result
				fmt.Println("Updated deployment!")
				fmt.Println(format.Backup(updated, cmd.RootArgs.Format))
			}
		},
	)
}
