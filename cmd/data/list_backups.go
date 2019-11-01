//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Brautigam Gergely
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
)

func init() {
	cmd.InitCommand(
		cmd.ListCmd,
		&cobra.Command{
			Use:   "backups",
			Short: "List backups",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				deploymentID string
				from         string
				to           string
			}{}
			f.StringVar(&cargs.deploymentID, "deployment-id", "", "The ID of the deployment to list backups for")
			f.StringVar(&cargs.from, "from", "", "Request backups that are created at or after this timestamp")
			f.StringVar(&cargs.to, "to", "", "Request backups that are created before this timestamp")
			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				deploymentID, argsUsed := cmd.OptOption("deployment-id", cargs.deploymentID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				backupc := backup.NewBackupServiceClient(conn)
				ctx := cmd.ContextWithToken()

				req := backup.ListBackupsRequest{
					DeploymentId: deploymentID,
				}

				if len(cargs.from) > 0 {
					from, err := time.Parse(time.RFC3339, cargs.from)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to parse from date. Accepted format is time.RFC3339")
					}
					req.From, err = types.TimestampProto(from)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to convert date to proto timestamp")
					}
				}

				if len(cargs.to) > 0 {
					to, err := time.Parse(time.RFC3339, cargs.from)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to parse to date. Accepted format is time.RFC3339")
					}
					req.To, err = types.TimestampProto(to)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to convert date to proto timestamp")
					}
				}

				// Fetch backups
				list, err := backupc.ListBackups(ctx, &req)
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to list backups")
				}

				// Show result
				fmt.Println(format.BackupList(list.Items, cmd.RootArgs.Format))
			}
		},
	)
}
