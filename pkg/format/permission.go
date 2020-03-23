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
// Author Ewout Prangsma
//

package format

import (
	"sort"
	"strings"

	"github.com/arangodb-managed/apis/common/auth"
)

type permissionRow struct {
	API   string
	Kind  string
	Verbs []string
}

// PermissionList returns a list of permissions formatted for humans.
func PermissionList(list []string, opts Options) string {
	var rows []permissionRow
	for _, p := range list {
		api, kind, verb, err := auth.ParsePermission(p)
		if err == nil {
			found := false
			for i, r := range rows {
				if r.API == api && r.Kind == kind {
					rows[i].Verbs = append(rows[i].Verbs, verb)
					found = true
					break
				}
			}
			if !found {
				rows = append(rows, permissionRow{API: api, Kind: kind, Verbs: []string{verb}})
			}
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		a, b := rows[i], rows[j]
		if a.API < b.API {
			return true
		}
		if a.Kind < b.Kind {
			return true
		}
		return false
	})

	return formatList(opts, rows, func(i int) []kv {
		x := rows[i]
		sort.Strings(x.Verbs)
		return []kv{
			kv{"api", x.API},
			kv{"kind", x.Kind},
			kv{"verbs", strings.Join(x.Verbs, ", ")},
		}
	}, false)
}
