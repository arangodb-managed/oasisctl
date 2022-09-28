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
		cmd.ListCmd,
		&cobra.Command{
			Use:   "notebookmodels",
			Short: "List notebook models",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				OrganizationID string
				ProjectID      string
				DeploymetnID   string
			}{}

			f.StringVarP(&cargs.OrganizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization that deployment is in")
			f.StringVarP(&cargs.ProjectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project that deployment is in")
			f.StringVarP(&cargs.DeploymetnID, "deployment-id", "d", "", "Identifier of the deployment that the notebook has to run next to")

			c.Run = func(c *cobra.Command, args []string) {
				log := cmd.CLILog

				deploymentID, argsUsed := cmd.ReqOption("deployment-id", cargs.DeploymetnID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				conn := cmd.MustDialAPI()
				notebookc := notebook.NewNotebookServiceClient(conn)
				datac := data.NewDataServiceClient(conn)
				rmc := rm.NewResourceManagerServiceClient(conn)
				ctx := cmd.ContextWithToken()

				selection.MustSelectDeployment(ctx, log, deploymentID, cargs.ProjectID, cargs.OrganizationID, datac, rmc)

				list, err := notebookc.ListNotebookModels(ctx, &notebook.ListNotebookModelsRequest{
					DeploymentId: deploymentID,
				})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to list notebook models")
				}
				fmt.Println(format.NotebookModelList(list.GetItems(), cmd.RootArgs.Format))
			}
		},
	)
}
