//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
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
		}
	}, false)
}
