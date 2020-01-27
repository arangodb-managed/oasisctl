//
// DISCLAIMER
//
// Copyright 2019 ArangoDB GmbH, Cologne, Germany
//
// Author Robert Stam
//

package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// ListServersCmd is based for other commands
	ListServersCmd = &cobra.Command{
		Use:   "servers",
		Short: "List servers information",
		Run:   ShowUsage,
	}
)

func init() {
	ListCmd.AddCommand(ListServersCmd)
}
