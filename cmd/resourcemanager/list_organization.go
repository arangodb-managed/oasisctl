//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package rm

import (
	"github.com/spf13/cobra"

	"github.com/arangodb-managed/oasisctl/cmd"
)

var (
	// listOrganizationCmd is root for various `list organization ...` commands
	listOrganizationCmd = &cobra.Command{
		Use:   "organization",
		Short: "List organization resources",
		Run:   cmd.ShowUsage,
	}
)

func init() {
	cmd.ListCmd.AddCommand(listOrganizationCmd)
}
