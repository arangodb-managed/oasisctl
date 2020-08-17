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

package cmd

import (
	"context"
	"crypto/tls"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/arangodb-managed/apis/common/auth"

	"github.com/arangodb-managed/oasisctl/pkg/format"
)

var (
	// RootCmd is the root (and only) command of this service
	RootCmd = &cobra.Command{
		Use:              "oasisctl",
		Short:            "ArangoDB Oasis",
		Long:             "ArangoDB Oasis. The Managed Cloud for ArangoDB",
		Run:              ShowUsage,
		PersistentPreRun: rootCmdPersistentPreRun,
	}

	CLILog = zerolog.New(zerolog.ConsoleWriter{
		Out:     os.Stderr,
		NoColor: !supportsColor(),
	}).With().Timestamp().Logger()
	RootArgs struct {
		Token    string
		endpoint string
		Format   format.Options
	}
)

const (
	// Prefix of all environment variables
	envKeyPrefix  = "OASIS_"
	envKeyProgram = envKeyPrefix + "EXECUTABLE"
	apiPortSuffix = ":443"
)

func init() {
	f := RootCmd.PersistentFlags()
	// Persistent flags
	defaultEndpoint := envOrDefault("ENDPOINT", "api.cloud.arangodb.com")
	f.StringVar(&RootArgs.Token, "token", "", "Token used to authenticate at ArangoDB Oasis")
	f.StringVar(&RootArgs.endpoint, "endpoint", defaultEndpoint, "API endpoint of the ArangoDB Oasis")
	f.StringVar(&RootArgs.Format.Format, "format", DefaultFormat(), "Output format (table|json)")
}

// ShowUsage shows usage of the given command on stdout.
func ShowUsage(cmd *cobra.Command, args []string) {
	cmd.Usage()
}

// Called before actual command run.
// This function is used to hide a default token (from environment variable)
// from the usage output.
func rootCmdPersistentPreRun(cmd *cobra.Command, args []string) {
	if RootArgs.Token == "" {
		RootArgs.Token = envOrDefault("TOKEN", "")
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
	if RootArgs.Token == "" {
		CLILog.Fatal().Msg("--token missing")
	}
	return auth.WithAccessToken(context.Background(), RootArgs.Token)
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

// InitCommand adds the given command to the given parent and called the flag initialization
// function.
func InitCommand(parent, cmd *cobra.Command, flagInit func(c *cobra.Command, f *flag.FlagSet)) *cobra.Command {
	if parent != nil {
		parent.AddCommand(cmd)
	}
	flagInit(cmd, cmd.Flags())
	return cmd
}
