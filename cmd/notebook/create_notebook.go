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

	data "github.com/arangodb-managed/apis/data/v1"
	notebook "github.com/arangodb-managed/apis/notebook/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
	"github.com/spf13/cobra"
)

func init() {
	cmd.InitCommand(
		cmd.CreateCmd,
		&cobra.Command{
			Use:   "notebook",
			Short: "Create a new notebook",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				OrganizationID string
				ProjectID      string
				DeploymetnID   string
				Name           string
				Description    string
				DiskSize       int32
				NotebookModel  string
			}{}

			f.StringVarP(&cargs.OrganizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization to create the notebook in")
			f.StringVarP(&cargs.ProjectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project to create the notebook in")
			f.StringVarP(&cargs.DeploymetnID, "deployment-id", "d", "", "Identifier of the deployment that the notebook has to run next to")
			f.StringVarP(&cargs.NotebookModel, "notebook-model", "m", "", "Identifier of the notebook model that the notebook has to use")
			f.StringVarP(&cargs.Name, "name", "n", "", "Name of the notebook")
			f.StringVarP(&cargs.Description, "description", "", "", "Description of the notebook")
			f.Int32VarP(&cargs.DiskSize, "disk-size", "s", 0, "Disk size in GiB that has to be attached to given notebook")

			c.Run = func(c *cobra.Command, args []string) {
				log := cmd.CLILog

				deploymentID, argsUsed := cmd.ReqOption("deployment-id", cargs.DeploymetnID, args, 0)
				name, _ := cmd.ReqOption("name", cargs.Name, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				conn := cmd.MustDialAPI()
				notebookc := notebook.NewNotebookServiceClient(conn)
				datac := data.NewDataServiceClient(conn)
				rmc := rm.NewResourceManagerServiceClient(conn)
				ctx := cmd.ContextWithToken()

				selection.MustSelectDeployment(ctx, log, deploymentID, cargs.ProjectID, cargs.OrganizationID, datac, rmc)

				n := notebook.Notebook{
					DeploymentId: deploymentID,
					Name:         name,
					Description:  cargs.Description,
					Model: &notebook.ModelSpec{
						NotebookModelId: cargs.NotebookModel,
						DiskSize:        cargs.DiskSize,
					},
				}
				created, err := notebookc.CreateNotebook(ctx, &n)
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to create notebook")
				}
				format.DisplaySuccess(cmd.RootArgs.Format)
				fmt.Println(format.Notebook(created, cmd.RootArgs.Format))
			}
		},
	)
}
