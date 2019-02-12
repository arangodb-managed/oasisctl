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
	// listOrganizationCmd is root for various `list organization ...` commands
	listOrganizationCmd = &cobra.Command{
		Use:   "organization",
		Short: "List organization resources",
		Run:   showUsage,
	}
)

func init() {
	listCmd.AddCommand(listOrganizationCmd)
}
