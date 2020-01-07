//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package iam

import (
	"github.com/spf13/cobra"

	"github.com/arangodb-managed/oasisctl/cmd"
)

var (
	// listGroupCmd is root for various `list group ...` commands
	listGroupCmd = &cobra.Command{
		Use:   "group",
		Short: "List group resources",
		Run:   cmd.ShowUsage,
	}
)

func init() {
	cmd.ListCmd.AddCommand(listGroupCmd)
}
