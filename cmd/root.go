//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package cmd

import (
	"context"
	"crypto/tls"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/arangodb-managed/adbcloud/pkg/format"
	"github.com/arangodb-managed/apis/common/auth"
)

var (
	// RootCmd is the root (and only) command of this service
	RootCmd = &cobra.Command{
		Use:   "adbcloud",
		Short: "ArangoDB Cloud",
		Run:   showUsage,
	}

	cliLog   = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	rootArgs struct {
		token    string
		endpoint string
		format   format.Options
	}
)

const (
	// Prefix of all environment variables
	envKeyPrefix  = "ADBCLOUD_"
	apiPortSuffix = ":443"
)

func init() {
	f := RootCmd.PersistentFlags()
	// Persistent flags
	defaultToken := envOrDefault("TOKEN", "")
	defaultEndpoint := envOrDefault("ENDPOINT", "cloud.adbtest.xyz")
	f.StringVar(&rootArgs.token, "token", defaultToken, "Token used to authenticate at ArangoDB Cloud")
	f.StringVar(&rootArgs.endpoint, "endpoint", defaultEndpoint, "API endpoint of the ArangoDB Cloud")
}

// Show usage of the given command
func showUsage(cmd *cobra.Command, args []string) {
	cmd.Usage()
}

// envOrDefault returns the value from an environment value with given key
// or if no such environment variable exists, the given default value.
func envOrDefault(envKeySuffix string, defaultValue string) string {
	if v := os.Getenv(envKeyPrefix + envKeySuffix); v != "" {
		return v
	}
	return defaultValue
}

// mustDialAPI dials the ArangoDB Cloud API
func mustDialAPI() *grpc.ClientConn {
	// Set up a connection to the server.
	tc := credentials.NewTLS(&tls.Config{})
	conn, err := grpc.Dial(rootArgs.endpoint+apiPortSuffix, grpc.WithTransportCredentials(tc))
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to connect to ArangoDB Cloud API")
	}
	return conn
}

// contextWithToken returns a context with access token in it.
func contextWithToken() context.Context {
	if rootArgs.token == "" {
		cliLog.Fatal().Msg("--token missing")
	}
	return auth.WithAccessToken(context.Background(), rootArgs.token)
}

// reqOption returns given value if not empty.
// Fails with clear error message when not set.
func reqOption(key, value string) string {
	if value == "" {
		cliLog.Fatal().Msgf("--%s missing", key)
	}
	return value
}
