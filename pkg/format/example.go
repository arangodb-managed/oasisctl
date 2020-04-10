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

package format

import (
	example "github.com/arangodb-managed/apis/example/v1"
)

// Example returns a single example dataset formatted for humans.
func Example(x *example.ExampleDataset, opts Options) string {
	return formatObject(opts,
		kv{"id", x.Id},
		kv{"name", x.Name},
		kv{"description", x.Description},
		kv{"url", x.Url},
		kv{"guide", x.Guide},
		kv{"created-at", formatTime(opts, x.CreatedAt)},
	)
}

// ExampleList returns a list of example datasets.
func ExampleList(list []*example.ExampleDataset, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			{"id", x.Id},
			{"name", x.Name},
			{"description", x.Description},
			{"url", x.Url},
			{"created-at", formatTime(opts, x.CreatedAt)},
		}
	}, false)
}
