//
// DISCLAIMER
//
// Copyright 2020-2024 ArangoDB GmbH, Cologne, Germany
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

package data

import (
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	backup "github.com/arangodb-managed/apis/backup/v1"
	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
)

var createBackupCmd = cmd.InitCommand(
	cmd.CreateCmd,
	&cobra.Command{
		Use:   "backup",
		Short: "Create backup ...",
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
					tp := timestamppb.New(t)
					if err := tp.CheckValid(); err != nil {
						log.Fatal().Err(err).Msg("Failed to convert from time to proto time")
					}
					b.AutoDeletedAt = tp
				}
			} else {
				if cargs.autoDeletedAt == 0 {
					cargs.autoDeletedAt = 6
				}
				t := time.Now().Add(time.Duration(cargs.autoDeletedAt) * time.Hour)
				tp := timestamppb.New(t)
				if err := tp.CheckValid(); err != nil {
					log.Fatal().Err(err).Msg("Failed to convert from time to proto time")
				}
				b.AutoDeletedAt = tp
			}

			result, err := backupc.CreateBackup(ctx, b)

			if err != nil {
				log.Fatal().Err(err).Msg("Failed to create backup")
			}

			// Show result
			format.DisplaySuccess(cmd.RootArgs.Format)
			fmt.Println(format.Backup(result, cmd.RootArgs.Format))
		}
	},
)
