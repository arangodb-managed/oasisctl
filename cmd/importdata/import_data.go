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
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	common "github.com/arangodb-managed/apis/common/v1"
	data "github.com/arangodb-managed/apis/data/v1"
	arangocopy "github.com/arangodb-managed/arangocopy/pkg"
	"github.com/arangodb-managed/oasisctl/cmd"
)

func init() {
	cmd.InitCommand(
		cmd.RootCmd,
		&cobra.Command{
			Use:   "import",
			Short: "Import data from a local database or from another remote database into an Oasis deployment.",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				source                  arangocopy.Connection
				destinationDeploymentID string
				includedDatabases       []string
				excludedDatabases       []string
				includedCollections     []string
				excludedCollections     []string
				includedViews           []string
				excludedViews           []string
				includedGraphs          []string
				excludedGraphs          []string
				force                   bool
				maxParallelCollections  int
				batchSize               int
				maxRetries              int
				queryTTL                time.Duration
				noProgressBar           bool
			}{}
			f.StringVar(&cargs.source.Address, "source-address", "", "Source database address to copy data from.")
			f.StringVar(&cargs.source.Username, "source-username", "", "Source database username if required.")
			f.StringVar(&cargs.source.Password, "source-password", "", "Source database password if required.")
			f.StringVarP(&cargs.destinationDeploymentID, "destination-deployment-id", "d", "", "Destination deployment id to import data into. It can be provided instead of address, username and password.")
			f.IntVarP(&cargs.maxParallelCollections, "maximum-parallel-collections", "m", 10, "Maximum number of collections being copied in parallel.")
			f.StringSliceVar(&cargs.includedDatabases, "included-database", []string{}, "A list of database names which should be included. If provided, only these databases will be copied.")
			f.StringSliceVar(&cargs.excludedDatabases, "excluded-database", []string{}, "A list of database names which should be excluded. Exclusion takes priority over inclusion.")
			f.StringSliceVar(&cargs.includedCollections, "included-collection", []string{}, "A list of collection names which should be included. If provided, only these collections will be copied.")
			f.StringSliceVar(&cargs.excludedCollections, "excluded-collection", []string{}, "A list of collections names which should be excluded. Exclusion takes priority over inclusion.")
			f.StringSliceVar(&cargs.includedViews, "included-view", []string{}, "A list of view names which should be included. If provided, only these views will be copied.")
			f.StringSliceVar(&cargs.excludedViews, "excluded-view", []string{}, "A list of view names which should be excluded. Exclusion takes priority over inclusion.")
			f.StringSliceVar(&cargs.includedGraphs, "included-graph", []string{}, "A list of graph names which should be included. If provided, only these graphs will be copied.")
			f.StringSliceVar(&cargs.excludedGraphs, "excluded-graph", []string{}, "A list of graph names which should be excluded. Exclusion takes priority over inclusion.")
			f.BoolVarP(&cargs.force, "force", "f", false, "Force the copy automatically overwriting everything at destination.")
			f.IntVarP(&cargs.batchSize, "batch-size", "b", 4096, "The number of documents to write at once.")
			f.IntVarP(&cargs.maxRetries, "max-retries", "r", 9, "The number of maximum retries attempts. Increasing this number will also increase the exponential fallback timer.")
			f.BoolVar(&cargs.noProgressBar, "no-progress-bar", false, "Disable the progress bar but still have partial progress output.")
			f.DurationVar(&cargs.queryTTL, "query-ttl", time.Hour*2, "Cursor TTL defined as a duration.")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				_, argsUsed := cmd.ReqOption("source-address", cargs.source.Address, args, 0)
				destinationDeploymentID, argsUsed := cmd.ReqOption("destination-deployment-id", cargs.destinationDeploymentID, args, 1)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				destination := arangocopy.Connection{}
				conn := cmd.MustDialAPI()
				datac := data.NewDataServiceClient(conn)
				ctx := cmd.ContextWithToken()
				depl, err := datac.GetDeployment(ctx, &common.IDOptions{Id: destinationDeploymentID})
				if err != nil {
					log.Fatal().Err(err).Str("deployment-id", destinationDeploymentID).Msg("Failed to get Deployment with id.")
				}
				creds, err := datac.GetDeploymentCredentials(ctx, &data.DeploymentCredentialsRequest{DeploymentId: destinationDeploymentID})
				if err != nil {
					log.Fatal().Err(err).Str("deployment-id", destinationDeploymentID).Msg("Failed to get Deployment credentials.")
				}
				destination.Address = depl.GetStatus().GetEndpoint()
				destination.Username = creds.Username
				destination.Password = creds.Password

				// Create copier
				copier, err := arangocopy.NewCopier(arangocopy.Config{
					Destination:                destination,
					Source:                     cargs.source,
					Force:                      cargs.force,
					BatchSize:                  cargs.batchSize,
					MaxRetries:                 cargs.maxRetries,
					QueryTTL:                   cargs.queryTTL,
					NoProgressBar:              cargs.noProgressBar,
					IncludedViews:              cargs.includedViews,
					ExcludedViews:              cargs.excludedViews,
					IncludedGraphs:             cargs.includedGraphs,
					ExcludedGraphs:             cargs.excludedGraphs,
					IncludedDatabases:          cargs.includedDatabases,
					ExcludedDatabases:          cargs.excludedDatabases,
					IncludedCollections:        cargs.includedCollections,
					ExcludedCollections:        cargs.excludedCollections,
					MaximumParallelCollections: cargs.maxParallelCollections,
				}, arangocopy.Dependencies{
					Logger:   log,
					Verifier: arangocopy.NewNoopVerifier(),
				})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to start copy operation.")
				}

				// Start copy operation
				if err := copier.Copy(); err != nil {
					log.Fatal().Err(err).Msg("Failed to copy. Please try again after the issue is resolved.")
				}
				log.Info().Msg("Success!")
			}
		},
	)
}
