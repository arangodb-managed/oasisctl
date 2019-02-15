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
	// GetCmd is root for various `get ...` commands
	GetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get information",
		Run:   ShowUsage,
	}
)

func init() {
	RootCmd.AddCommand(GetCmd)
}
