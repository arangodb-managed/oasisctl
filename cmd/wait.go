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
	// WaitCmd is root for various `wait ...` commands
	WaitCmd = &cobra.Command{
		Use:   "wait",
		Short: "Wait for a status change",
		Run:   ShowUsage,
	}
)

func init() {
	RootCmd.AddCommand(WaitCmd)
}
