//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Author Gergely Brautigam
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
