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
	// rejectOrganizationCmd is root for various `reject organization ...` commands
	rejectOrganizationCmd = &cobra.Command{
		Use:   "organization",
		Short: "Reject organization related invites",
		Run:   showUsage,
	}
)

func init() {
	rejectCmd.AddCommand(rejectOrganizationCmd)
}
