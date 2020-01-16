//
// DISCLAIMER
//
// Copyright 2020 ArangoDB Inc, Cologne, Germany
//
// Author Gergely Brautigam
//

package iam

import (
	"github.com/spf13/cobra"

	"github.com/arangodb-managed/oasisctl/cmd"
)

var (
	// addMembersCmd is root for various `add group ...` commands
	addMembersCmd = &cobra.Command{
		Use:   "group",
		Short: "Add group resources",
		Run:   cmd.ShowUsage,
	}
)

func init() {
	cmd.AddCmd.AddCommand(addMembersCmd)
}
