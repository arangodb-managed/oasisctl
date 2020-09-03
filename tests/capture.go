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
	"bytes"
	"io"
	"os"
)

// CaptureOutput captures os.Stdout on which oasisctl is logging its output.
func CaptureOutput(f func() error) (string, error) {
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = stdout }()

	if err := f(); err != nil {
		return "", err
	}

	if err := w.Close(); err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return "", err
	}
	return buf.String(), nil
}
