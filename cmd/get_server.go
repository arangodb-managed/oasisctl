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
	// GetServerCmd is based for other commands
	GetServerCmd = &cobra.Command{
		Use:   "server",
		Short: "Get server information",
		Run:   ShowUsage,
	}
)

func init() {
	GetCmd.AddCommand(GetServerCmd)
}
