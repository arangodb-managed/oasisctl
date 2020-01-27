//
// DISCLAIMER
//
// Copyright 2019 ArangoDB GmbH, Cologne, Germany
//
// Author Ewout Prangsma
//

package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// UpdateCmd is root for various `update ...` commands
	UpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update resources",
		Run:   ShowUsage,
	}
)

func init() {
	RootCmd.AddCommand(UpdateCmd)
}
