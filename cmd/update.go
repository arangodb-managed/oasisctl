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
	// updateCmd is root for various `update ...` commands
	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update resources",
		Run:   showUsage,
	}
)

func init() {
	RootCmd.AddCommand(updateCmd)
}
