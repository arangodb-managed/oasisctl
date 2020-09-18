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

	"github.com/arangodb-managed/oasisctl/pkg/selection"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	backup "github.com/arangodb-managed/apis/backup/v1"
	common "github.com/arangodb-managed/apis/common/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
)

func init() {
	cmd.InitCommand(
		deleteBackupCmd,
		&cobra.Command{
			Use:   "policy",
			Short: "Delete a backup policy for a given ID.",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				id             string
				organizationID string
				projectID      string
			}{}
			f.StringVarP(&cargs.id, "id", "i", "", "Identifier of the backup policy")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				id, argsUsed := cmd.OptOption("id", cargs.id, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				backupc := backup.NewBackupServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Fetch backup policy
				item := selection.MustSelectBackupPolicy(ctx, log, id, backupc)

				// Delete backup policy
				if _, err := backupc.DeleteBackupPolicy(ctx, &common.IDOptions{Id: item.GetId()}); err != nil {
					log.Fatal().Err(err).Msg("Failed to delete backup policy")
				}

				// Show result
				fmt.Println("Deleted backup policy!")
			}
		},
	)
}
