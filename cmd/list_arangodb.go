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
	// ListArangoDBCmd is based for other commands
	ListArangoDBCmd = &cobra.Command{
		Use:   "arangodb",
		Short: "List ArangoDB information",
		Run:   ShowUsage,
	}
)

func init() {
	ListCmd.AddCommand(ListArangoDBCmd)
}
