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
	// deleteCmd is root for various `delete ...` commands
	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete resources",
		Run:   showUsage,
	}
)

func init() {
	RootCmd.AddCommand(deleteCmd)
}
