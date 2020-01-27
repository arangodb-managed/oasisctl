//
// DISCLAIMER
//
// Copyright 2019 ArangoDB GmbH, Cologne, Germany
//
// Author Ewout Prangsma
//

package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// BackupCmd is root for various `backup ...` commands
	BackupCmd = &cobra.Command{
		Use:   "backup",
		Short: "Backup commands",
		Run:   ShowUsage,
	}
)

func init() {
	RootCmd.AddCommand(BackupCmd)
}
