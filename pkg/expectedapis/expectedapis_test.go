//
// DISCLAIMER
//
// Copyright 2026 ArangoDB GmbH, Cologne, Germany
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

package expectedapis

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func apisJSONPath(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve test file path")
	}
	return filepath.Join(filepath.Dir(filename), "..", "..", "apis.json")
}

func TestExpectedAPIVersionsMatchApisJSON(t *testing.T) {
	content, err := os.ReadFile(apisJSONPath(t))
	if err != nil {
		t.Fatalf("read apis.json: %v", err)
	}
	var committed map[string]string
	if err := json.Unmarshal(content, &committed); err != nil {
		t.Fatalf("unmarshal apis.json: %v", err)
	}

	generated := ExpectedAPIVersions()
	if len(generated) != len(committed) {
		t.Fatalf("api count mismatch: registry has %d apis, apis.json has %d — run make update-apis-json", len(generated), len(committed))
	}
	for id, want := range committed {
		got, ok := generated[id]
		if !ok {
			t.Errorf("apis.json contains %q but registry does not — run make update-apis-json", id)
			continue
		}
		if got != want {
			t.Errorf("version mismatch for %q: registry=%q apis.json=%q — run make update-apis-json", id, got, want)
		}
	}
	for id := range generated {
		if _, ok := committed[id]; !ok {
			t.Errorf("registry contains %q but apis.json does not — run make update-apis-json", id)
		}
	}
}
