//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Robert Stam
//

package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// ListNodeSizesCmd is based for other commands
	ListNodeSizesCmd = &cobra.Command{
		Use:   "node-sizes",
		Short: "List node sizes information",
		Run:   ShowUsage,
	}
)

func init() {
	ListCmd.AddCommand(ListNodeSizesCmd)
}
