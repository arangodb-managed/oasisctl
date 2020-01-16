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
	// AddCmd is root for various `add ...` commands
	AddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add resources",
		Run:   ShowUsage,
	}
)

func init() {
	RootCmd.AddCommand(AddCmd)
}
