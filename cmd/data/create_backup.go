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
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	backup "github.com/arangodb-managed/apis/backup/v1"
	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
)

func init() {
	cmd.InitCommand(
		cmd.CreateCmd,
		&cobra.Command{
			Use:   "backup",
			Short: "Create a new backup",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				name          string
				deploymentID  string
				policyID      string
				description   string
				autoDeletedAt int
				upload        bool
			}{}
			f.StringVar(&cargs.name, "name", "", "Name of the deployment")
			f.StringVar(&cargs.deploymentID, "deployment-id", "", "ID of the deployment")
			f.StringVar(&cargs.description, "description", "", "Description of the backup")
			f.BoolVar(&cargs.upload, "upload", false, "The backup should be uploaded")
			f.IntVar(&cargs.autoDeletedAt, "auto-deleted-at", 0, "Time (h) until auto delete of the backup")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				name, argsUsed := cmd.ReqOption("name", cargs.name, args, 0)
				deploymentID, argsUsed := cmd.ReqOption("deployment-id", cargs.deploymentID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				backupc := backup.NewBackupServiceClient(conn)
				ctx := cmd.ContextWithToken()

				b := &backup.Backup{
					Name:         name,
					DeploymentId: deploymentID,
					Description:  cargs.description,
				}

				if cargs.upload {
					b.Upload = true
					if cargs.autoDeletedAt != 0 {
						t := time.Now().Add(time.Duration(cargs.autoDeletedAt) * time.Hour)
						tp, err := types.TimestampProto(t)
						if err != nil {
							log.Fatal().Err(err).Msg("Failed to convert from time to proto time")
						}
						b.AutoDeletedAt = tp
					}
				} else {
					if cargs.autoDeletedAt == 0 {
						cargs.autoDeletedAt = 6
					}
					t := time.Now().Add(time.Duration(cargs.autoDeletedAt) * time.Hour)
					tp, err := types.TimestampProto(t)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to convert from time to proto time")
					}
					b.AutoDeletedAt = tp
				}

				result, err := backupc.CreateBackup(ctx, b)

				if err != nil {
					log.Fatal().Err(err).Msg("Failed to create backup")
				}

				// Show result
				fmt.Println("Success!")
				fmt.Println(format.Backup(result, cmd.RootArgs.Format))
			}
		},
	)
}
