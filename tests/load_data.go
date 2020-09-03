//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
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

// LoadArguments loads arguments from a defined argument file.
func LoadArguments(f string) ([]string, error) {
	data, err := loadOutput(fmt.Sprintf("%s.args", f))
	if err != nil {
		return nil, err
	}
	// if there are no delimiters, just return the content
	if !strings.Contains(string(data), " ") {
		return []string{string(data)}, nil
	}

	args := strings.Split(string(data), " ")
	if len(args) == 0 {
		return nil, errors.New("empty argument list")
	}
	return args, nil
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
