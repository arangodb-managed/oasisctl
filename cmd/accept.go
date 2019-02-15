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
	// AcceptCmd is root for various `accept ...` commands
	AcceptCmd = &cobra.Command{
		Use:   "accept",
		Short: "Accept invites",
		Run:   ShowUsage,
	}
)

func init() {
	RootCmd.AddCommand(AcceptCmd)
}
