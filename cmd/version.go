//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package cmd

import (
	"fmt"

	"github.com/coreos/go-semver/semver"
	"github.com/spf13/cobra"

	"github.com/arangodb-managed/oasisctl/pkg/format"
)

var (
	currentVersion *semver.Version
	// versionCmd is the command to show the current version of this tool
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show the current version of this tool",
		Run:   runVersionCmd,
	}
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

// SetVersion must be called at bootstrap to pass the current build version
func SetVersion(releaseVersion string) {
	currentVersion = semver.Must(semver.NewVersion(releaseVersion))
}

// Run the service
func runVersionCmd(cmd *cobra.Command, args []string) {
	fmt.Println(format.CLIVersion(currentVersion.String(), RootArgs.Format))
}
