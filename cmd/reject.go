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
	// rejectCmd is root for various `reject ...` commands
	rejectCmd = &cobra.Command{
		Use:   "reject",
		Short: "Reject invites",
		Run:   showUsage,
	}
)

func init() {
	RootCmd.AddCommand(rejectCmd)
}
