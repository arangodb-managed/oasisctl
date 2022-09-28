//
// DISCLAIMER
//
// Copyright 2022 ArangoDB GmbH, Cologne, Germany
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

package notebook

import (
	"fmt"

	flag "github.com/spf13/pflag"

	common "github.com/arangodb-managed/apis/common/v1"
	notebook "github.com/arangodb-managed/apis/notebook/v1"
	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/spf13/cobra"
)

func init() {
	cmd.InitCommand(
		cmd.GetCmd,
		&cobra.Command{
			Use:   "notebook",
			Short: "Get notebook",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				OrganizationID string
				ProjectID      string
				ID             string
			}{}

			f.StringVarP(&cargs.ID, "notebook-id", "n", "", "Identifier of the notebook")

			c.Run = func(c *cobra.Command, args []string) {
				log := cmd.CLILog

				id, argsUsed := cmd.ReqOption("notebook-id", cargs.ID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				conn := cmd.MustDialAPI()
				notebookc := notebook.NewNotebookServiceClient(conn)
				ctx := cmd.ContextWithToken()

				notebook, err := notebookc.GetNotebook(ctx, &common.IDOptions{
					Id: id,
				})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to get notebook")
				}

				fmt.Println(format.Notebook(notebook, cmd.RootArgs.Format))
			}
		},
	)
}
