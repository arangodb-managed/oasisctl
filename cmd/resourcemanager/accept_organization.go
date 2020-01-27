//
// DISCLAIMER
//
// Copyright 2019 ArangoDB GmbH, Cologne, Germany
//
// Author Ewout Prangsma
//

package rm

import (
	"github.com/spf13/cobra"

	"github.com/arangodb-managed/oasisctl/cmd"
)

var (
	// acceptOrganizationCmd is root for various `accept organization ...` commands
	acceptOrganizationCmd = &cobra.Command{
		Use:   "organization",
		Short: "Accept organization related invites",
		Run:   cmd.ShowUsage,
	}
)

func init() {
	cmd.AcceptCmd.AddCommand(acceptOrganizationCmd)
}
