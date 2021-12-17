//
// DISCLAIMER
//
// Copyright 2020-2021 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
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
	
. <(oasisctl completion [bash|fish|powershell|zsh])
	
To configure your bash shell to load completions for each session add to your bashrc
	
# ~/.bashrc or ~/.profile
. <(oasisctl completion bash)
`,
		ValidArgs: []string{"", "bash", "fish", "powershell", "zsh"},
		Run: func(cmd *cobra.Command, args []string) {
			shell := "bash"
			if len(args) > 0 {
				shell = args[0]
			}
			switch shell {
			case "bash":
				RootCmd.GenBashCompletion(os.Stdout)
			case "fish":
				RootCmd.GenFishCompletion(os.Stdout, true)
			case "zsh":
				RootCmd.GenZshCompletion(os.Stdout)
			case "powershell":
				RootCmd.GenPowerShellCompletion(os.Stdout)
			default:
				CLILog.Fatal().Str("shell", shell).Msg("Unknown shell")
			}
		},
	}
)

func init() {
	RootCmd.AddCommand(completionCmd)
}
