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

	"github.com/arangodb-managed/apis/common/auth"
	"github.com/arangodb-managed/oasis/pkg/format"
)

var (
	// RootCmd is the root (and only) command of this service
	RootCmd = &cobra.Command{
		Use:              "oasis",
		Short:            "ArangoDB Oasis",
		Long:             "ArangoDB Oasis. The Managed Cloud for ArangoDB",
		Run:              ShowUsage,
		PersistentPreRun: rootCmdPersistentPreRun,
	}

	CLILog   = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	RootArgs struct {
		token    string
		endpoint string
		Format   format.Options
	}
)

const (
	// Prefix of all environment variables
	envKeyPrefix  = "OASIS_"
	apiPortSuffix = ":443"
)

func init() {
	f := RootCmd.PersistentFlags()
	// Persistent flags
	defaultEndpoint := envOrDefault("ENDPOINT", "cloud.adbtest.xyz")
	f.StringVar(&RootArgs.token, "token", "", "Token used to authenticate at ArangoDB Oasis")
	f.StringVar(&RootArgs.endpoint, "endpoint", defaultEndpoint, "API endpoint of the ArangoDB Oasis")
}

// Show usage of the given command
func ShowUsage(cmd *cobra.Command, args []string) {
	cmd.Usage()
}

// Called before actual command run.
// This function is used to hide a default token (from environment variable)
// from the usage output.
func rootCmdPersistentPreRun(cmd *cobra.Command, args []string) {
	if RootArgs.token == "" {
		RootArgs.token = envOrDefault("TOKEN", "")
	}
}

// envOrDefault returns the value from an environment value with given key
// or if no such environment variable exists, the given default value.
func envOrDefault(envKeySuffix string, defaultValue string) string {
	if v := os.Getenv(envKeyPrefix + envKeySuffix); v != "" {
		return v
	}
	return defaultValue
}

// MustDialAPI dials the ArangoDB Oasis API
func MustDialAPI() *grpc.ClientConn {
	// Set up a connection to the server.
	tc := credentials.NewTLS(&tls.Config{})
	conn, err := grpc.Dial(RootArgs.endpoint+apiPortSuffix, grpc.WithTransportCredentials(tc))
	if err != nil {
		CLILog.Fatal().Err(err).Msg("Failed to connect to ArangoDB Oasis API")
	}
	return conn
}

// ContextWithToken returns a context with access token in it.
func ContextWithToken() context.Context {
	if RootArgs.token == "" {
		CLILog.Fatal().Msg("--token missing")
	}
	return auth.WithAccessToken(context.Background(), RootArgs.token)
}

// ReqOption returns given value if not empty.
// Fails with clear error message when not set.
// Returns: option-value, number-of-args-used(0|argIndex+1)
func ReqOption(key, value string, args []string, argIndex int) (string, int) {
	if value != "" {
		return value, 0
	}
	if len(args) > argIndex {
		return args[argIndex], argIndex + 1
	}
	CLILog.Fatal().Msgf("--%s missing", key)
	return "", 0
}

// OptOption returns given value if not empty.
// Returns: option-value, number-of-args-used(0|argIndex+1)
func OptOption(key, value string, args []string, argIndex int) (string, int) {
	if value != "" {
		return value, 0
	}
	if len(args) > argIndex {
		return args[argIndex], argIndex + 1
	}
	return "", 0
}

// MustCheckNumberOfArgs compares the number of arguments with the expected
// number of arguments.
// If there is a difference a fatal error is raised.
func MustCheckNumberOfArgs(args []string, expectedNumberOfArgs int) {
	if len(args) > expectedNumberOfArgs {
		CLILog.Fatal().Msg("Too many arguments")
	}
	if len(args) < expectedNumberOfArgs {
		CLILog.Fatal().Msg("Too few arguments")
	}
}
