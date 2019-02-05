//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// getCmd is root for various `get ...` commands
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get information",
		Run:   showUsage,
	}
)

func init() {
	RootCmd.AddCommand(getCmd)
}
