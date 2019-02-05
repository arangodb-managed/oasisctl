//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package cmd

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
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
		token string
	}
)

func init() {
	f := RootCmd.PersistentFlags()
	// Persistent flags
	f.StringVar(&rootArgs.token, "token", os.Getenv("ADBCLOUD_TOKEN"), "Token used to authenticate at ArangoDB Cloud")
}

// Show usage of the given command
func showUsage(cmd *cobra.Command, args []string) {
	cmd.Usage()
}
