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
	// acceptCmd is root for various `accept ...` commands
	acceptCmd = &cobra.Command{
		Use:   "accept",
		Short: "Accept invites",
		Run:   showUsage,
	}
)

func init() {
	RootCmd.AddCommand(acceptCmd)
}
