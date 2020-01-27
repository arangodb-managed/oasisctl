//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Author Gergely Brautigam
//

package iam

import (
	"github.com/spf13/cobra"

	"github.com/arangodb-managed/oasisctl/cmd"
)

var (
	// addGroupCmd is root for various `add group ...` commands
	addGroupCmd = &cobra.Command{
		Use:   "group",
		Short: "Add group resources",
		Run:   cmd.ShowUsage,
	}
)

func init() {
	cmd.AddCmd.AddCommand(addGroupCmd)
}
