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
	// rejectOrganizationCmd is root for various `reject organization ...` commands
	rejectOrganizationCmd = &cobra.Command{
		Use:   "organization",
		Short: "Reject organization related invites",
		Run:   cmd.ShowUsage,
	}
)

func init() {
	cmd.RejectCmd.AddCommand(rejectOrganizationCmd)
}
