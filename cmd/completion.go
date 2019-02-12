//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	// completionCmd generates a shell command line completion script
	completionCmd = &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long: `To load completion run
	
	. <(oasis completion)
	
	To configure your bash shell to load completions for each session add to your bashrc
	
	# ~/.bashrc or ~/.profile
	. <(oasis completion)
	`,
		Run: func(cmd *cobra.Command, args []string) {
			RootCmd.GenBashCompletion(os.Stdout)
		},
	}
)

func init() {
	RootCmd.AddCommand(completionCmd)
}
