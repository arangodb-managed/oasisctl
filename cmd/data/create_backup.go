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
				name        string
				deploymenID string
				policyID    string
				// TODO add other fields
			}{}
			f.StringVar(&cargs.name, "name", "", "Name of the deployment")
			f.StringVar(&cargs.deploymenID, "deployment-id", "", "ID of the deployment")
			f.StringVar(&cargs.policyID, "policy-id", "", "ID of the backup policy to use")

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

				result, err := backupc.CreateBackup(ctx, &backup.Backup{
					Name:           name,
					DeploymentId:   deploymenID,
					BackupPolicyId: cargs.policyID,
					// TODO: AutoDeleteAt
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
