//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
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
// Author Brautigam Gergely
//

package data

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	backup "github.com/arangodb-managed/apis/backup/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/util"
)

var listBackupsCmd = cmd.InitCommand(
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
				var err error
				req.From, err = util.ParseTime(cargs.from)
				if err != nil {
					log.Fatal().Err(err)
				}
			}

			if len(cargs.to) > 0 {
				var err error
				req.To, err = util.ParseTime(cargs.to)
				if err != nil {
					log.Fatal().Err(err)
				}
			}

			// Fetch all backups
			var backups []*backup.Backup
			if err := backup.ForEachBackup(ctx, func(ctx context.Context, req *backup.ListBackupsRequest) (*backup.BackupList, error) {
				return backupc.ListBackups(ctx, req)
			}, req, func(ctx context.Context, backup *backup.Backup) error {
				backups = append(backups, backup)
				return nil
			}); err != nil {
				log.Fatal().Err(err).Msg("Failed to list backups")
			}

			// Show result
			fmt.Println(format.BackupList(backups, cmd.RootArgs.Format))
		}
	},
)
