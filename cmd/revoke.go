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
	// RevokeCmd is root for various `revoke ...` commands
	RevokeCmd = &cobra.Command{
		Use:   "revoke",
		Short: "Revoke keys & tokens",
		Run:   ShowUsage,
	}
)

func init() {
	RootCmd.AddCommand(RevokeCmd)
}
