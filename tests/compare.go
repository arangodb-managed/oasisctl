//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Author Gergely Brautigam
//

package tests

import (
	"fmt"
	"regexp"
)

// CompareResults regex compare the result with the output so keys like .*,
// can be used to check random ids, or timestamps or values which should match something but
// there is no possible way to define what that something will be.
func CompareResults(output []byte, regex []byte) bool {
	match := regexp.MustCompile(string(regex))
	if !match.Match(output) {
		fmt.Printf("Output: \n%s\n", string(output))
		fmt.Printf("Regex: \n%s\n", string(regex))
		return false
	}
	return true
}
