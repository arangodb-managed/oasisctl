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
	// ListEffectiveCmd is based for other commands
	ListEffectiveCmd = &cobra.Command{
		Use:   "effective",
		Short: "List effective information",
		Run:   ShowUsage,
	}
)

func init() {
	ListCmd.AddCommand(ListEffectiveCmd)
}
