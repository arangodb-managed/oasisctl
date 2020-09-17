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
// Author Gergely Brautigam
//

package data

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	backup "github.com/arangodb-managed/apis/backup/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.UnlockCmd,
		&cobra.Command{
			Use:   "policy",
			Short: "Unlock a backup policy",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				id string
			}{}
			f.StringVarP(&cargs.id, "backup-policy-id", "d", "", "Identifier of the backup policy")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				id, argsUsed := cmd.OptOption("backup-policy-id", cargs.id, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				backupc := backup.NewBackupServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Select a backup policy to update
				item := selection.MustSelectBackupPolicy(ctx, log, id, backupc)
				item.Locked = false

				updated, err := backupc.UpdateBackupPolicy(ctx, item)
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to unlock backup")
				}

				// Show result
				fmt.Println("Unlocked backup policy!")
				fmt.Println(format.BackupPolicy(updated, cmd.RootArgs.Format))
			}
		},
	)
}
