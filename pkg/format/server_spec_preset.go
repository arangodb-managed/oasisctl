//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Robert Stam
//

package format

import (
	"fmt"

	data "github.com/arangodb-managed/apis/data/v1"
)

// ServersSpecPresetList returns a list of servers spec presets formatted for humans.
func ServersSpecPresetList(list []*data.ServersSpecPreset, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			kv{"name", x.GetName()},
			kv{"default", formatBool(x.GetIsDefault())},
			kv{"coordinators", x.GetServers().GetCoordinators()},
			kv{"coordinator-memory-size", fmt.Sprintf("%d%s", x.GetServers().GetCoordinatorMemorySize(), "GB")},
			kv{"dbservers", x.GetServers().GetDbservers()},
			kv{"dbserver-memory-size", fmt.Sprintf("%d%s", x.GetServers().GetDbserverMemorySize(), "GB")},
			kv{"dbserver-disk-size", fmt.Sprintf("%d%s", x.GetServers().GetDbserverDiskSize(), "GB")},
		}
	}, false)
}
