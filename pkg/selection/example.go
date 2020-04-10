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
// Author Ewout Prangsma
//

package selection

import (
	"context"

	"github.com/rs/zerolog"

	common "github.com/arangodb-managed/apis/common/v1"
	example "github.com/arangodb-managed/apis/example/v1"
)

// MustSelectExampleDataset fetches an example dataset with given ID, name, or URL and fails if no such item is found.
func MustSelectExampleDataset(ctx context.Context, log zerolog.Logger, id string, examplec example.ExampleDatasetServiceClient) *example.ExampleDataset {
	result, err := SelectExampleDataset(ctx, log, id, examplec)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to select example dataset")
	}
	return result
}

// SelectExampleDataset fetches an example dataset with given ID, name, or URL or returns an error if not found.
func SelectExampleDataset(ctx context.Context, log zerolog.Logger, id string, examplec example.ExampleDatasetServiceClient) (*example.ExampleDataset, error) {
	result, err := examplec.GetExampleDataset(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) {
			// Try to lookup example dataset by name or URL
			list, err := examplec.ListExampleDatasets(ctx, &example.ListExampleDatasetsRequest{})
			if err == nil {
				for _, x := range list.Items {
					if x.GetName() == id || x.GetUrl() == id {
						return x, nil
					}
				}
			}
		}
		log.Debug().Err(err).Str("example_dataset", id).Msg("Failed to get example dataset")
		return nil, err
	}
	return result, nil
}
