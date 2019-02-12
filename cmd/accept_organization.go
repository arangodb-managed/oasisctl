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
	// acceptOrganizationCmd is root for various `accept organization ...` commands
	acceptOrganizationCmd = &cobra.Command{
		Use:   "organization",
		Short: "Accept organization related invites",
		Run:   showUsage,
	}
)

func init() {
	acceptCmd.AddCommand(acceptOrganizationCmd)
}
