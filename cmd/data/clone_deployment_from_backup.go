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

	replication "github.com/arangodb-managed/apis/replication/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.CloneDeploymentCmd,
		&cobra.Command{
			Use:   "backup",
			Short: "Clone a deployment from a backup.",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				backupID       string
				organizationID string
				regionID       string
				acceptTAndC    bool
			}{}
			f.StringVarP(&cargs.backupID, "backup-id", "b", "", "Clone a deployment from a backup using the backup's ID.")
			f.StringVarP(&cargs.regionID, "region-id", "r", "", "An optionally defined region in which the new deployment should be created in.")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization to create the clone in")
			f.BoolVar(&cargs.acceptTAndC, "accept", false, "Accept the current terms and conditions.")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				backupID, argsUsed := cmd.OptOption("backup-id", cargs.backupID, args, 0)
				regionID, argsUsed := cmd.OptOption("region-id", cargs.regionID, args, argsUsed)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				ctx := cmd.ContextWithToken()
				repl := replication.NewReplicationServiceClient(conn)

				req := &replication.CloneDeploymentFromBackupRequest{
					BackupId: backupID,
					RegionId: regionID,
				}
				if cargs.acceptTAndC {
					rmc := rm.NewResourceManagerServiceClient(conn)
					tandc := selection.MustSelectTermsAndConditions(ctx, log, "", cargs.organizationID, rmc)
					req.AcceptedTermsAndConditionsId = tandc.GetId()
				}

				// Clone deployment
				created, err := repl.CloneDeploymentFromBackup(ctx, req)
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to clone deployment")
				}

				// Show result
				format.DisplaySuccess(cmd.RootArgs.Format)
				fmt.Println(format.Deployment(created, nil, cmd.RootArgs.Format, false))
			}
		},
	)
}
