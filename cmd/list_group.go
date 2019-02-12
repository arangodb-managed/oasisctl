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
	// listGroupCmd is root for various `list group ...` commands
	listGroupCmd = &cobra.Command{
		Use:   "group",
		Short: "List group resources",
		Run:   showUsage,
	}
)

func init() {
	listCmd.AddCommand(listGroupCmd)
}
