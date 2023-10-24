//
// DISCLAIMER
//
// Copyright 2020-2023 ArangoDB GmbH, Cologne, Germany
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

package format

import (
	"fmt"
	"strings"

	data "github.com/arangodb-managed/apis/data/v1"
)

// NodeSizeList returns a list of node sizes.
func NodeSizeList(list []*data.NodeSize, cpuList []*data.CPUSize, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			{"id", x.GetId()},
			{"name", x.GetName()},
			{"max-disk-size", fmt.Sprintf("%d%s", x.GetMaxDiskSize(), "GiB")},
			{"min-disk-size", fmt.Sprintf("%d%s", x.GetMinDiskSize(), "GiB")},
			{"allowed-disk-sizes", formatAllowedDiskSizes(x.GetDiskSizes())},
			{"memory-size", fmt.Sprintf("%d%s", x.GetMemorySize(), "GiB")},
			{"cpu-size", formatCPUSize(x.GetCpuSize(), cpuList)},
		}
	}, true)
}

func formatAllowedDiskSizes(list []int32) string {
	if len(list) == 0 {
		return "any"
	}
	result := make([]string, 0, len(list))
	for _, x := range list {
		result = append(result, fmt.Sprintf("%d%s", x, "GiB"))
	}
	return strings.Join(result, ", ")
}

func formatCPUSize(id string, list []*data.CPUSize) string {
	for _, x := range list {
		if x.GetId() == id {
			return x.GetName()
		}
	}
	return id
}
