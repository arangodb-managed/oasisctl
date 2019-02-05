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
	// createCmd is root for various `create ...` commands
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create resources",
		Run:   showUsage,
	}
)

func init() {
	RootCmd.AddCommand(createCmd)
}
