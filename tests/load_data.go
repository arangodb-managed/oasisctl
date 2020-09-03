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
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

const argumentPrefix = "%"

// LoadCommands loads arguments from a defined argument file.
func LoadCommands(f string, args map[string]string) ([][]string, error) {
	data, err := loadOutput(fmt.Sprintf("%s.cmd", f))
	if err != nil {
		return nil, err
	}

	commands := strings.Split(string(data), "\n")

	if len(commands) == 0 {
		return nil, errors.New("empty commands")
	}

	// replace all the possible command argument placeholders.
	var replace []string
	for k, v := range args {
		replace = append(replace, argumentPrefix+k, v)
	}
	r := strings.NewReplacer(replace...)
	cmds := [][]string{}
	for _, l := range commands {
		l = r.Replace(l)
		cmd := strings.Split(l, " ")
		if len(cmd) == 0 {
			return nil, errors.New("empty command arguments")
		}
		cmds = append(cmds, cmd)
	}
	return cmds, nil
}

// LoadGoldenOutput loads the test's golden output.
func LoadGoldenOutput(f string) ([]byte, error) {
	return loadOutput(fmt.Sprintf("%s.golden.output", f))
}

// LoadErrorOutput loads the test's error output.
func LoadErrorOutput(f string) ([]byte, error) {
	return loadOutput(fmt.Sprintf("%s.error.output", f))
}

func loadOutput(f string) ([]byte, error) {
	return ioutil.ReadFile(fmt.Sprintf("./testdata/%s", f))
}
