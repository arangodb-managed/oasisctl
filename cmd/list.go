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
	// listCmd is root for various `list ...` commands
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List resources",
		Run:   showUsage,
	}
)

func init() {
	RootCmd.AddCommand(listCmd)
}
