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

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	backup "github.com/arangodb-managed/apis/backup/v1"
	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/format"
	"github.com/gogo/protobuf/types"
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
				deploymenID   string
				policyID      string
				description   string
				autoDeletedAt int
			}{}
			f.StringVar(&cargs.name, "name", "", "Name of the deployment")
			f.StringVar(&cargs.deploymenID, "deployment-id", "", "ID of the deployment")
			f.StringVar(&cargs.description, "description", "", "Description of the backup")
			f.IntVar(&cargs.autoDeletedAt, "autodeletedat", 6, "Time (h) until auto delete of the backup")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				name, argsUsed := cmd.ReqOption("name", cargs.name, args, 0)
				deploymenID, argsUsed := cmd.ReqOption("deployment-id", cargs.deploymenID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				backupc := backup.NewBackupServiceClient(conn)
				ctx := cmd.ContextWithToken()

				t := time.Now().Add(time.Duration(cargs.autoDeletedAt) * time.Hour)
				tp, err := types.TimestampProto(t)
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to convert from time to proto time")
				}

				result, err := backupc.CreateBackup(ctx, &backup.Backup{
					Name:          name,
					DeploymentId:  deploymenID,
					Description:   cargs.description,
					AutoDeletedAt: tp,
				})

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
