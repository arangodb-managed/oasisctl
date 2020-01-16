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
	// deleteMembersCmd is root for various `delete group ...` commands
	deleteMembersCmd = &cobra.Command{
		Use:   "group",
		Short: "Delete group resources",
		Run:   cmd.ShowUsage,
	}
)

func init() {
	cmd.DeleteCmd.AddCommand(deleteMembersCmd)
}
