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
	// updatePolicyCmd is root for various `update policy ...` commands
	updatePolicyCmd = &cobra.Command{
		Use:   "policy",
		Short: "Update a policy",
		Run:   cmd.ShowUsage,
	}
	// updatePolicyAddCmd is root for various `update policy add ...` commands
	updatePolicyAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add to a policy",
		Run:   cmd.ShowUsage,
	}
	// updatePolicyDeleteCmd is root for various `update policy delete ...` commands
	updatePolicyDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete from a policy",
		Run:   cmd.ShowUsage,
	}
)

func init() {
	cmd.UpdateCmd.AddCommand(updatePolicyCmd)
	updatePolicyCmd.AddCommand(updatePolicyAddCmd)
	updatePolicyCmd.AddCommand(updatePolicyDeleteCmd)
}
