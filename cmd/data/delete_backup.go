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
	common "github.com/arangodb-managed/apis/common/v1"

	"github.com/arangodb-managed/oasis/cmd"
)

func init() {
	cmd.InitCommand(
		cmd.DeleteCmd,
		&cobra.Command{
			Use:   "backup",
			Short: "Delete a backup for a given ID.",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				backupID string
			}{}
			f.StringVarP(&cargs.backupID, "id", "i", "", "Identifier of the backup")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				backupID, argsUsed := cmd.OptOption("id", cargs.backupID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				backupc := backup.NewBackupServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Delete backup
				if _, err := backupc.DeleteBackup(ctx, &common.IDOptions{Id: backupID}); err != nil {
					log.Fatal().Err(err).Msg("Failed to delete deployment")
				}

				// Show result
				fmt.Println("Deleted backup!")
			}
		},
	)
}
