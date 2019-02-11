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
	// getEffectiveCmd is based for other commands
	getEffectiveCmd = &cobra.Command{
		Use:   "effective",
		Short: "Get information",
		Run:   showUsage,
	}
)

func init() {
	getCmd.AddCommand(getEffectiveCmd)
}
