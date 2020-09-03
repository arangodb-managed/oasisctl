//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
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
// Author Gergely Brautigam
//

package tests

import (
	"testing"

	"github.com/arangodb-managed/oasisctl/cmd"
)

// RunCommands takes a *testing.T and a name of the test.
// Performs the steps of the test running the specified
// set of commands and checks the end results against the
// golden output.
func RunCommands(t *testing.T, name string, args map[string]string) {
	cmds, err := LoadCommands(name, args)
	if err != nil {
		t.Fatal(err)
	}

	// Run all commands and capture each output and finally, compare to the end result.
	var out string
	for _, c := range cmds {
		root := cmd.RootCmd
		root.SetArgs(c)
		o, err := CaptureOutput(root.Execute)
		if err != nil {
			t.Fatal(err)
		}
		out += o
	}

	regex, err := LoadGoldenOutput(name)
	if err != nil {
		t.Fatal(err)
	}

	if !CompareResults([]byte(out), regex) {
		t.FailNow()
	}
}
