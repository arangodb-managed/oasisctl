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

// copyDatabase creates a database at the destination.
func (c *copier) copyDatabase(ctx context.Context, db driver.Database) error {
	return c.backoffCall(ctx, func() error {
		if err := c.ensureDestinationDatabase(ctx, db.Name()); err != nil {
			return err
		}
		return nil
	})
}

// ensureDestinationDatabase ensures that a database exists at the destination.
func (c *copier) ensureDestinationDatabase(ctx context.Context, dbName string) error {
	c.Logger.Debug().Str("database-name", dbName).Msg("Ensuring database exists")
	if exists, err := c.destinationClient.DatabaseExists(ctx, dbName); err != nil {
		c.Logger.Warn().Err(err).Msg("Failed to get if database exists.")
		return err
	} else if exists {
		return nil
	}
	_, err := c.destinationClient.CreateDatabase(ctx, dbName, nil)
	return err
}
