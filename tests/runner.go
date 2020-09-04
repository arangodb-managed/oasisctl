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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

const (
	oasisctl = "oasisctl"
	windows  = "windows"
	linux    = "linux"
	darwin   = "darwin"
)

// RunCommands takes a *testing.T and a name of the test.
// Performs the steps of the test running the specified
// set of commands and checks the end results against the
// golden output.
func RunCommands(t *testing.T, compare string, args []string, fail bool) {
	// Find oasisctl executable
	cmd, err := lookupOasisctl()
	if err != nil {
		t.Fatal(err)
	}
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil && !fail {
		fmt.Println("Failed output: ", string(out))
		t.Fatal(err)
	}
	if !CompareResults(out, []byte(compare)) {
		t.FailNow()
	}
}

func lookupOasisctl() (string, error) {
	var path string
	op := runtime.GOOS
	arch := runtime.GOARCH
	switch op {
	case windows:
		path = filepath.Join("..", "..", "bin", windows, arch, oasisctl)
	case darwin:
		path = filepath.Join("..", "..", "bin", darwin, arch, oasisctl)
	case linux:
		path = filepath.Join("..", "..", "bin", linux, arch, oasisctl)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", err
	}
	return path, nil
}
