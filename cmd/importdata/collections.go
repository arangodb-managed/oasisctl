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
	"errors"
	"sort"

	"github.com/arangodb/go-driver"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

// copyCollections copies all collections for a database.
func (c *copier) copyCollections(ctx context.Context, db driver.Database) error {
	log := c.Logger
	log.Info().Msg("Beginning to copy over collection data.")
	var (
		collections   []driver.Collection
		destinationDB driver.Database
	)
	if err := c.backoffCall(ctx, func() error {
		colls, err := db.Collections(ctx)
		if err != nil {
			c.Logger.Error().Err(err).Msg("Failed to list collections for source database.")
			return err
		}
		collections = colls
		return nil
	}); err != nil {
		return err
	}
	if err := c.backoffCall(ctx, func() error {
		ddb, err := c.destinationClient.Database(ctx, db.Name())
		if err != nil {
			c.Logger.Error().Err(err).Msg("Failed to get destination database.")
			return nil
		}
		destinationDB = ddb
		return nil
	}); err != nil {
		return err
	}

	collections = c.filterCollections(collections)

	// gather all props into a map to minimise the sorting function and not repeate the props
	// call later when copying the data.
	propsMap := make(map[string]driver.CollectionProperties)
	for _, coll := range collections {
		if err := c.backoffCall(ctx, func() error {
			prop, err := coll.Properties(ctx)
			if err != nil {
				c.Logger.Error().Err(err).Msg("Failed to get properties.")
				return err
			}
			propsMap[coll.Name()] = prop
			return nil
		}); err != nil {
			return err
		}
	}
	if err := c.sortCollections(collections, propsMap); err != nil {
		log.Error().Err(err).Msg("Failed to sort collections")
		return err
	}

	readCtx := driver.WithQueryStream(ctx, true)
	readCtx = driver.WithQueryBatchSize(readCtx, c.BatchSize)
	restoreCtx := driver.WithIsRestore(ctx, true)
	var g errgroup.Group
	sem := semaphore.NewWeighted(int64(c.MaximumParallelCollections))
	if c.Dependencies.Spinner != nil {
		c.Dependencies.Spinner.Start()
	}

	// Create the collections sequentially here
	for _, sourceColl := range collections {
		if err := c.createCollection(restoreCtx, destinationDB, sourceColl, propsMap[sourceColl.Name()]); err != nil {
			c.Logger.Error().Err(err).Msg("Failed to ensure destination collection.")
			return err
		}
	}

	// Start the data copy operation
	for _, sourceColl := range collections {
		sourceColl := sourceColl
		g.Go(func() error {
			// Ensure semaphore.
			if err := sem.Acquire(ctx, 1); err != nil {
				return err
			}
			defer sem.Release(1)

			props, ok := propsMap[sourceColl.Name()]
			if !ok {
				return errors.New("no properties found for collection")
			}
			if props.IsSystem {
				// skip system collections
				return nil
			}

			var destinationColl driver.Collection
			if err := c.backoffCall(ctx, func() error {
				dColl, err := destinationDB.Collection(ctx, sourceColl.Name())
				if err != nil {
					c.Logger.Error().Err(err).Msg("Failed to ensure destination collection.")
					return err
				}
				destinationColl = dColl
				return nil
			}); err != nil {
				return err
			}

			// Copy over all indexes for this collection.
			if err := c.copyIndexes(restoreCtx, sourceColl, destinationColl); err != nil {
				c.Logger.Error().Err(err).Str("collection", sourceColl.Name()).Msg("Failed to copy all indexes.")
				return err
			}

			bindVars := map[string]interface{}{
				"@c": sourceColl.Name(),
			}
			var cursor driver.Cursor
			if err := c.backoffCall(ctx, func() error {
				cr, err := db.Query(readCtx, "FOR d IN @@c RETURN d", bindVars)
				if err != nil {
					c.Logger.Error().Err(err).Str("collection", sourceColl.Name()).Msg("Failed to query source database for collection.")
					return err
				}
				cursor = cr
				return nil
			}); err != nil {
				return err
			}
			defer cursor.Close()
			batch := make([]interface{}, 0, c.BatchSize)
			for {
				var (
					d           interface{}
					noMoreError bool
				)
				if err := c.backoffCall(ctx, func() error {
					if _, err := cursor.ReadDocument(readCtx, &d); driver.IsNoMoreDocuments(err) {
						noMoreError = true
						return nil
					} else if err != nil {
						c.Logger.Error().Err(err).Str("collection", sourceColl.Name()).Msg("Read documents failed.")
						return err
					}
					batch = append(batch, d)
					return nil
				}); err != nil {
					return err
				}

				if (noMoreError && len(batch) > 0) || len(batch) >= c.BatchSize {
					if err := c.backoffCall(ctx, func() error {
						if _, _, err := destinationColl.CreateDocuments(restoreCtx, batch); err != nil {
							c.Logger.Error().Err(err).Str("collection", sourceColl.Name()).Interface("document", d).Msg("Creating a document failed.")
							return err
						}
						batch = make([]interface{}, 0, c.BatchSize)
						return nil
					}); err != nil {
						return err
					}
				}

				if noMoreError {
					break
				}
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		log.Error().Err(err).Msg("One of the workers failed to copy data.")
		return err
	}
	if c.Dependencies.Spinner != nil {
		c.Dependencies.Spinner.Stop()
	}
	log.Debug().Str("source-database", db.Name()).Msg("Done copying database data.")
	return nil
}

// copyIndexes copies all indexes for a collection to destination collection.
func (c *copier) copyIndexes(ctx context.Context, sourceColl driver.Collection, destinationColl driver.Collection) error {
	var indexes []driver.Index
	if err := c.backoffCall(ctx, func() error {
		idxs, err := sourceColl.Indexes(ctx)
		if err != nil {
			return err
		}
		indexes = idxs
		return nil
	}); err != nil {
		return err
	}
	for _, index := range indexes {
		if err := c.backoffCall(ctx, func() error {
			switch index.Type() {
			case driver.TTLIndex:
				var field string
				if len(index.Fields()) > 0 {
					field = index.Fields()[0]
				}
				if _, _, err := destinationColl.EnsureTTLIndex(ctx, field, index.ExpireAfter(), &driver.EnsureTTLIndexOptions{
					Name: index.UserName(),
				}); err != nil {
					return err
				}
			case driver.PersistentIndex:
				if _, _, err := destinationColl.EnsurePersistentIndex(ctx, index.Fields(), &driver.EnsurePersistentIndexOptions{
					Name:   index.UserName(),
					Unique: index.Unique(),
					Sparse: index.Sparse(),
				}); err != nil {
					return err
				}
			case driver.SkipListIndex:
				if _, _, err := destinationColl.EnsureSkipListIndex(ctx, index.Fields(), &driver.EnsureSkipListIndexOptions{
					Name:          index.UserName(),
					Unique:        index.Unique(),
					Sparse:        index.Sparse(),
					NoDeduplicate: !index.Deduplicate(),
				}); err != nil {
					return err
				}
			case driver.HashIndex:
				if _, _, err := destinationColl.EnsureHashIndex(ctx, index.Fields(), &driver.EnsureHashIndexOptions{
					Name:          index.UserName(),
					Unique:        index.Unique(),
					Sparse:        index.Sparse(),
					NoDeduplicate: !index.Deduplicate(),
				}); err != nil {
					return err
				}
			case driver.FullTextIndex:
				if _, _, err := destinationColl.EnsureFullTextIndex(ctx, index.Fields(), &driver.EnsureFullTextIndexOptions{
					Name:      index.UserName(),
					MinLength: index.MinLength(),
				}); err != nil {
					return err
				}
			case driver.GeoIndex:
				if _, _, err := destinationColl.EnsureGeoIndex(ctx, index.Fields(), &driver.EnsureGeoIndexOptions{
					Name:    index.UserName(),
					GeoJSON: index.GeoJSON(),
				}); err != nil {
					return err
				}
			case driver.EdgeIndex:
				// These are automatically created.
			case driver.PrimaryIndex:
				// These are automatically created.
			default:
				return errors.New("unknown index type " + string(index.Type()))
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

// createCollection creates a collection on the destination database.
func (c *copier) createCollection(ctx context.Context, db driver.Database, coll driver.Collection, props driver.CollectionProperties) error {
	if props.IsSystem {
		// skip system collections
		return nil
	}
	exists, err := db.CollectionExists(ctx, coll.Name())
	if err != nil {
		c.Logger.Warn().Err(err).Msg("Failed to get if collection exists.")
		return err
	}
	if exists {
		return nil
	}
	options := &driver.CreateCollectionOptions{
		JournalSize:       int(props.JournalSize),
		ReplicationFactor: props.ReplicationFactor,
		WriteConcern:      props.WriteConcern,
		WaitForSync:       props.WaitForSync,
		DoCompact:         &props.DoCompact,
		CacheEnabled:      &props.CacheEnabled,
		ShardKeys:         props.ShardKeys,
		NumberOfShards:    props.NumberOfShards,
		IsSystem:          false,
		Type:              props.Type,
		KeyOptions: &driver.CollectionKeyOptions{
			AllowUserKeys: props.KeyOptions.AllowUserKeys,
			Type:          props.KeyOptions.Type,
		},
		DistributeShardsLike: props.DistributeShardsLike,
		IsSmart:              false,
		ShardingStrategy:     props.ShardingStrategy,
	}
	if props.SmartJoinAttribute != "" {
		options.IsSmart = true
		options.SmartJoinAttribute = props.SmartJoinAttribute
	}
	if err := c.backoffCall(ctx, func() error {
		if _, err := db.CreateCollection(ctx, coll.Name(), options); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// sortCollections sorts a collection list on the following premis:
// No DistributedShardsLike comes before DistributedShardsLike. In case two collections don't have this setting
// a vertex collection comes before an edge collection.
func (c *copier) sortCollections(collections []driver.Collection, m map[string]driver.CollectionProperties) error {
	sort.SliceStable(collections, func(i, j int) bool {
		pi := m[collections[i].Name()]
		pj := m[collections[j].Name()]

		// Collections without DistributeShardsLike must come before collections which have that setting
		if pi.DistributeShardsLike == "" && pj.DistributeShardsLike != "" {
			return true
		} else if pi.DistributeShardsLike != "" && pj.DistributeShardsLike == "" {
			return false
		}

		// Vertex collections should come before edge collections
		if pi.Type == driver.CollectionTypeDocument && pj.Type == driver.CollectionTypeEdge {
			return true
		} else if pi.Type == driver.CollectionTypeEdge && pj.Type == driver.CollectionTypeDocument {
			return false
		}

		// Lastely, sort by name
		return pi.Name < pj.Name
	})
	return nil
}
