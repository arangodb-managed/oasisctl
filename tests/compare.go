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
	"regexp"
)

// CompareOutput regex compare the result with the output so keys like .*,
// can be used to check random ids, or timestamps or values which should match something but
// there is no possible way to define what that something will be.
func CompareOutput(output []byte, regex []byte) bool {
	regex = escapeRegexSpecificCharacters(regex)
	match := regexp.MustCompile(string(regex))
	if !match.Match(output) {
		fmt.Println("Output: ")
		verbose(string(output))
		fmt.Println()
		fmt.Println("Regex: ")
		verbose(string(regex))
		return false
	}
	return true
}

// escapeRegexSpecificCharacters will escape the most common interfering characters from the regex.
func escapeRegexSpecificCharacters(s []byte) []byte {
	var output []byte
	for _, c := range s {
		switch c {
		case '[', ']', '|':
			output = append(output, '\\', c)
		default:
			output = append(output, c)
		}
	}
	return output
}

// verbose makes the output more verbose so number of characters can be seen like space, newline, tabs.
func verbose(s string) {
	for _, r := range s {
		if r == ' ' {
			fmt.Print(".")
		} else if r == '\t' {
			fmt.Print("\\t")
			fmt.Print("\t")
		} else if r == '\n' {
			fmt.Print("\\n")
			fmt.Println()
		} else {
			fmt.Print(string(r))
		}
	}
}
