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
	"fmt"

	data "github.com/arangodb-managed/apis/data/v1"
)

// NodeSizeList returns a list of node sizes.
func NodeSizeList(list []*data.NodeSize, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			kv{"name", x.Name},
			kv{"max-disk-size", fmt.Sprintf("%d%s", x.MaxDiskSize, "GB")},
			kv{"min-disk-size", fmt.Sprintf("%d%s", x.MinDiskSize, "GB")},
			kv{"memory-size", fmt.Sprintf("%d%s", x.MemorySize, "GB")},
			kv{"cpu-size", x.CpuSize},
		}
	}, false)
}
