//
// DISCLAIMER
//
// Copyright 2022-2026 ArangoDB GmbH, Cologne, Germany
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

package notebook

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/cobra"

	"github.com/arangodb-managed/oasisctl/cmd"
)

// OAS-12992: old ArangoGraphML Notebooks CLI is retired. Keep command names so
// existing scripts get a clear failure instead of a missing-command surprise.
// This is unrelated to ArangoDB Platform "Data Science Platform".
const deprecatedNotebookMsg = "notebook commands are deprecated and no longer available (old ArangoGraphML/Notebooks stack retired)"

func init() {
	registerDeprecated(cmd.CreateCmd, "notebook", "Create a new notebook (deprecated)")
	registerDeprecated(cmd.GetCmd, "notebook", "Get a notebook (deprecated)")
	registerDeprecated(cmd.ListCmd, "notebooks", "List notebooks (deprecated)")
	registerDeprecated(cmd.ListCmd, "notebookmodels", "List notebook models (deprecated)")
	registerDeprecated(cmd.UpdateCmd, "notebook", "Update notebook (deprecated)")
	registerDeprecated(cmd.DeleteCmd, "notebook", "Delete a notebook (deprecated)")
	registerDeprecated(cmd.PauseCmd, "notebook", "Pause a notebook (deprecated)")
	registerDeprecated(cmd.ResumeCmd, "notebook", "Resume a notebook (deprecated)")
}

func registerDeprecated(parent *cobra.Command, use, short string) {
	cmd.InitCommand(
		parent,
		&cobra.Command{
			Use:        use,
			Short:      short,
			Deprecated: deprecatedNotebookMsg,
			// Ignore leftover flags (--notebook-id, etc.) so Run always emits the deprecation error.
			FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		},
		func(c *cobra.Command, _ *flag.FlagSet) {
			c.Run = func(*cobra.Command, []string) {
				cmd.CLILog.Fatal().Msg(deprecatedNotebookMsg)
			}
		},
	)
}
