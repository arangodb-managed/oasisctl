//
// DISCLAIMER
//
// Copyright 2019 ArangoDB GmbH, Cologne, Germany
//
// Author Gergely Brautigam
//

package data

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	backup "github.com/arangodb-managed/apis/backup/v1"
	v1 "github.com/arangodb-managed/apis/common/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
)

func init() {
	cmd.InitCommand(
		cmd.BackupCmd,
		&cobra.Command{
			Use:   "download",
			Short: "Download a backup",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				ID string
			}{}
			f.StringVarP(&cargs.ID, "id", "i", "", "Identifier of the backup")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				id, argsUsed := cmd.OptOption("id", cargs.ID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				backupc := backup.NewBackupServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Fetch backup
				_, err := backupc.DownloadBackup(ctx, &v1.IDOptions{Id: id})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to fetch backup")
				}

				// Show result
				fmt.Println("Backup download started successfully!")
			}
		},
	)
}
