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
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/briandowns/spinner"
	"github.com/cenkalti/backoff"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh/terminal"
)

// Copier copies database content from source address to destination address
type Copier interface {
	// Copy copies over every database, collection, view, etc, from source to destination.
	// At the destination everything will be overwritten.
	Copy() error
}

// Connection encapsulates connection details for a database.
type Connection struct {
	Address  string
	Username string
	Password string
}

// Config defines configuration for this copier.
type Config struct {
	// Source database connection details
	Source Connection
	// Destination database connection details
	Destination Connection
	// Destination Deployment ID. Optionally define a deployment ID as a destination.
	DeploymentID string
	// A list of database names to be included in the copy operation. If define, only these names will be selected.
	IncludedDatabases []string
	// A list of database names to be excluded from the copy operation.
	ExcludedDatabases []string
	// A list of collection names to be included in the copy operation. If define, only these names will be selected.
	IncludedCollections []string
	// A list of collection names to be excluded from the copy operation.
	ExcludedCollections []string
	// A list of view names to be included in the copy operation. If define, only these names will be selected.
	IncludedViews []string
	// A list of view names to be excluded from the copy operation.
	ExcludedViews []string
	// A list of graph names to be included in the copy operation. If define, only these names will be selected.
	IncludedGraphs []string
	// A list of graph names to be excluded from the copy operation.
	ExcludedGraphs []string
	// Forces the copy operation ignoring the confirm dialog.
	Force bool
	// Number of parallel collection copies underway.
	MaximumParallelCollections int
	// The batch size of the cursor.
	BatchSize int
	// MaxRetries defines the number of retries the backoff will do.
	MaxRetries int
}

// Dependencies defines dependencies for the copier.
type Dependencies struct {
	Logger  zerolog.Logger
	Spinner *spinner.Spinner
}

type copier struct {
	Config
	Dependencies
	sourceClient      driver.Client
	destinationClient driver.Client
	databaseInclude   map[string]struct{}
	databaseExclude   map[string]struct{}
	collectionInclude map[string]struct{}
	collectionExclude map[string]struct{}
	viewInclude       map[string]struct{}
	viewExclude       map[string]struct{}
	graphInclude      map[string]struct{}
	graphExclude      map[string]struct{}
}

// NewCopier returns a new copier with given a given set of configurations.
func NewCopier(cfg Config, deps Dependencies) (Copier, error) {
	c := &copier{
		Config:       cfg,
		Dependencies: deps,
	}
	// Set up source client.
	if client, err := c.getClient("Source", cfg.Source.Address, cfg.Source.Username, cfg.Source.Password); err != nil {
		c.Logger.Error().Err(err).Msg("Failed to connect to source address.")
		return nil, err
	} else {
		c.sourceClient = client
	}
	// Set up destination client.
	if client, err := c.getClient("Destination", cfg.Destination.Address, cfg.Destination.Username, cfg.Destination.Password); err != nil {
		c.Logger.Error().Err(err).Msg("Failed to connect to destination address.")
		return nil, err
	} else {
		c.destinationClient = client
	}
	// Set up spinner
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		c.Dependencies.Spinner = spinner.New(spinner.CharSets[34], 100*time.Millisecond)
	}

	// Set up filters
	c.databaseInclude = setupMap(c.Config.IncludedDatabases)
	c.databaseExclude = setupMap(c.Config.ExcludedDatabases)
	c.collectionInclude = setupMap(c.Config.IncludedCollections)
	c.collectionExclude = setupMap(c.Config.ExcludedCollections)
	c.viewInclude = setupMap(c.Config.IncludedViews)
	c.viewExclude = setupMap(c.Config.ExcludedViews)
	c.graphInclude = setupMap(c.Config.IncludedGraphs)
	c.graphExclude = setupMap(c.Config.ExcludedGraphs)
	return c, nil
}

// A small helper to setup a map for filters.
func setupMap(data []string) map[string]struct{} {
	m := make(map[string]struct{})
	for _, f := range data {
		m[f] = struct{}{}
	}
	return m
}

// getClient creates a client pointing to address and tests if that connection works.
func (c *copier) getClient(prefix, address, username, password string) (driver.Client, error) {
	log := c.Logger
	// Open a connection
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{address},
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to open conncetion to address")
		return nil, err
	}
	// Create the client
	cfg := driver.ClientConfig{
		Connection: conn,
	}
	if username != "" && password != "" {
		cfg.Authentication = driver.BasicAuthentication(username, password)
	}
	client, err := driver.NewClient(cfg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create driver client")
		return nil, err
	}
	// Test a connection to the database
	version, err := client.Version(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		return nil, err
	}
	log.Info().Msgf("%s: Version at address (%s) is %s", prefix, address, version.String())
	return client, nil
}

// Copy copies over every database, collection, view, etc, from source to destination.
// At the destination everything will be overwritten.
func (c *copier) Copy() error {
	log := c.Logger

	if ok, err := c.displayConfirmation(); err != nil {
		return err
	} else if !ok {
		log.Info().Msg("Cancelling operation.")
		return nil
	}

	ctx := context.Background()
	// Gather all databases
	var databases []driver.Database
	if err := c.backoffCall(ctx, func() error {
		dbs, err := c.sourceClient.Databases(ctx)
		if err != nil {
			c.Logger.Error().Err(err).Msg("Failed to get databases for source.")
			return err
		}
		databases = dbs
		return nil
	}); err != nil {
		return err
	}

	databases = c.filterDatabases(databases)

	for _, db := range databases {
		if err := c.copyDatabase(ctx, db); err != nil {
			return err
		}
		log.Info().Msg("Done with databases.")
		if err := c.copyCollections(ctx, db); err != nil {
			return err
		}
		log.Info().Msg("Done with collections.")
		if err := c.copyViews(ctx, db); err != nil {
			return err
		}
		log.Info().Msg("Done with viewes.")
		if err := c.copyGraphs(ctx, db); err != nil {
			return err
		}
		log.Info().Msg("Done with graphs.")
	}

	return nil
}

// displayConfirmation will only display the confirm question if the terminal is an interactive one.
// otherwise, will fail.
func (c *copier) displayConfirmation() (bool, error) {
	if terminal.IsTerminal(int(os.Stdout.Fd())) && !c.Force {
		var response string
		fmt.Print("Please confirm copy operation (y/N) ")
		fmt.Scanln(&response)
		if response != "y" {
			log.Info().Msg("Halting operation.")
			return false, nil
		}
		return true, nil
	} else if c.Force {
		return true, nil
	}
	return false, errors.New("either use an interactive terminal or define --force flag")
}

// backoffCall is a convenient wrapper around backoff Retry.
func (c *copier) backoffCall(ctx context.Context, f func() error) error {
	if err := backoff.Retry(f, backoff.WithContext(backoff.WithMaxRetries(backoff.NewExponentialBackOff(), uint64(c.MaxRetries)), ctx)); err != nil {
		c.Logger.Error().Err(err).Msg("Backoff eventually failed.")
		return err
	}
	return nil
}
