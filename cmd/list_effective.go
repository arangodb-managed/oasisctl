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
	// listEffectiveCmd is based for other commands
	listEffectiveCmd = &cobra.Command{
		Use:   "effective",
		Short: "List effective information",
		Run:   showUsage,
	}
)

func init() {
	listCmd.AddCommand(listEffectiveCmd)
}
