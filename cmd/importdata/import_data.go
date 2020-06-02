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

package importdata

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	common "github.com/arangodb-managed/apis/common/v1"
	data "github.com/arangodb-managed/apis/data/v1"
	"github.com/arangodb-managed/oasisctl/cmd"
)

func init() {
	cmd.InitCommand(
		cmd.RootCmd,
		&cobra.Command{
			Use:   "import",
			Short: "Import data from a local database or from another remote database into a deployment by ID or by endpoint credentials.",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				source                 Connection
				destination            Connection
				deploymentID           string
				includedDatabases      []string
				excludedDatabases      []string
				includedCollections    []string
				excludedCollections    []string
				includedViews          []string
				excludedViews          []string
				includedGraphs         []string
				excludedGraphs         []string
				force                  bool
				maxParallelCollections int
				batchSize              int
				maxRetries             int
			}{}
			f.StringVarP(&cargs.source.Address, "source-address", "s", "", "Source database address to copy data from.")
			f.StringVar(&cargs.source.Username, "source-username", "", "Source database username if required.")
			f.StringVar(&cargs.source.Password, "source-password", "", "Source database password if required.")
			f.StringVar(&cargs.destination.Address, "destination-address", "", "Destination database address to copy data to.")
			f.StringVar(&cargs.destination.Username, "destination-username", "", "Destination database username if required.")
			f.StringVar(&cargs.destination.Password, "destination-password", "", "Destination database password if required.")
			f.StringVarP(&cargs.deploymentID, "destination-deployment-id", "d", "", "Destination deployment id to import data into. It can be provided instead of address, username and password.")
			f.IntVarP(&cargs.maxParallelCollections, "maximum-parallel-collections", "m", 10, "Maximum number of collections being copied in parallel.")
			f.StringSliceVar(&cargs.includedDatabases, "included-database", []string{}, "A list of database names which should be included. If provided, only these databases will be copied.")
			f.StringSliceVar(&cargs.excludedDatabases, "exluded-database", []string{}, "A list of database names which should be excluded. Exclusion takes priority over inclusion.")
			f.StringSliceVar(&cargs.includedCollections, "included-collection", []string{}, "A list of collection names which should be included. If provided, only these collections will be copied.")
			f.StringSliceVar(&cargs.excludedCollections, "excluded-collection", []string{}, "A list of collections names which should be excluded. Exclusion takes priority over inclusion.")
			f.StringSliceVar(&cargs.includedViews, "included-view", []string{}, "A list of view names which should be included. If provided, only these views will be copied.")
			f.StringSliceVar(&cargs.excludedViews, "excluded-view", []string{}, "A list of view names which should be excluded. Exclusion takes priority over inclusion.")
			f.StringSliceVar(&cargs.includedGraphs, "included-graph", []string{}, "A list of graph names which should be included. If provided, only these graphs will be copied.")
			f.StringSliceVar(&cargs.excludedGraphs, "excluded-graph", []string{}, "A list of graph names which should be excluded. Exclusion takes priority over inclusion.")
			f.BoolVarP(&cargs.force, "force", "f", false, "Force the copy automatically overwriting everything at destination.")
			f.IntVarP(&cargs.batchSize, "batch-size", "b", 4096, "The number of documents to write at once.")
			f.IntVarP(&cargs.maxRetries, "max-retries", "r", 9, "The number of maximum retries attempts. Increasing this number will also increase the exponential fallback timer.")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				_, argsUsed := cmd.ReqOption("source-address", cargs.source.Address, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				destination := cargs.destination

				if cargs.deploymentID != "" {
					conn := cmd.MustDialAPI()
					datac := data.NewDataServiceClient(conn)
					ctx := cmd.ContextWithToken()
					depl, err := datac.GetDeployment(ctx, &common.IDOptions{Id: cargs.deploymentID})
					if err != nil {
						log.Fatal().Err(err).Str("deployment-id", cargs.deploymentID).Msg("Failed to get Deployment with id.")
					}
					destination.Address = depl.GetStatus().GetEndpoint()
					creds, err := datac.GetDeploymentCredentials(ctx, &data.DeploymentCredentialsRequest{DeploymentId: cargs.deploymentID})
					if err != nil {
						log.Fatal().Err(err).Str("deployment-id", cargs.deploymentID).Msg("Failed to get Deployment credentials.")
					}
					destination.Password = creds.Password
					destination.Username = creds.Username
				}
				// Create copier
				copier, err := NewCopier(Config{
					Force:                      cargs.force,
					Source:                     cargs.source,
					BatchSize:                  cargs.batchSize,
					MaxRetries:                 cargs.maxRetries,
					Destination:                destination,
					IncludedViews:              cargs.includedViews,
					ExcludedViews:              cargs.excludedViews,
					IncludedGraphs:             cargs.includedGraphs,
					ExcludedGraphs:             cargs.excludedGraphs,
					IncludedDatabases:          cargs.includedDatabases,
					ExcludedDatabases:          cargs.excludedDatabases,
					IncludedCollections:        cargs.includedCollections,
					ExcludedCollections:        cargs.excludedCollections,
					MaximumParallelCollections: cargs.maxParallelCollections,
				}, Dependencies{
					Logger: log,
				})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to start copy operation.")
				}

				// Start copy operatio
				if err := copier.Copy(); err != nil {
					log.Fatal().Err(err).Msg("Failed to copy. Please try again after the issue is resolved.")
				}
				log.Info().Msg("Success!")
			}
		},
	)
}
