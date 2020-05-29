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
	"context"

	"github.com/arangodb/go-driver"
)

// copyGraphs copies all graphs from source database to destination database.
func (c *copier) copyGraphs(ctx context.Context, source driver.Database) error {
	log := c.Logger
	ctx = driver.WithIsRestore(ctx, true)
	var destination driver.Database
	if err := c.backoffCall(ctx, func() error {
		// Get the destination database
		destDB, err := c.destinationClient.Database(ctx, source.Name())
		if err != nil {
			log.Error().Err(err).Msg("Failed to get destination database")
			return err
		}
		destination = destDB
		return nil
	}); err != nil {
		return err
	}

	var graphs []driver.Graph
	if err := c.backoffCall(ctx, func() error {
		gs, err := source.Graphs(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to list source graphs")
			return err
		}
		graphs = gs
		return nil
	}); err != nil {
		return err
	}
	graphs = c.filterGraphs(graphs)

	// Get the replication factor of the target system.
	var destinationReplicationFactor int
	if err := c.backoffCall(ctx, func() error {
		info, err := destination.Info(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get destination database info.")
			return err
		}
		destinationReplicationFactor = info.ReplicationFactor
		return nil
	}); err != nil {
		return err
	}

	for _, g := range graphs {
		var (
			exists bool
		)
		if err := c.backoffCall(ctx, func() error {
			if ok, err := destination.GraphExists(ctx, g.Name()); err != nil {
				log.Error().Err(err).Msg("Error checking if graph exists.")
				return err
			} else {
				exists = ok
			}
			return nil
		}); err != nil {
			return err
		}

		if exists {
			continue
		}

		if err := c.backoffCall(ctx, func() error {
			replFactor := g.ReplicationFactor()
			if replFactor < destinationReplicationFactor {
				replFactor = destinationReplicationFactor
			}
			if _, err := destination.CreateGraph(ctx, g.Name(), &driver.CreateGraphOptions{
				OrphanVertexCollections: g.OrphanCollections(),
				EdgeDefinitions:         g.EdgeDefinitions(),
				IsSmart:                 g.IsSmart(),
				SmartGraphAttribute:     g.SmartGraphAttribute(),
				NumberOfShards:          g.NumberOfShards(),
				ReplicationFactor:       replFactor,
				WriteConcern:            g.WriteConcern(),
			}); err != nil {
				log.Error().Err(err).Msg("Failed to create graph.")
				return err
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}
