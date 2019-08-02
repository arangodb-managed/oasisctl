//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package format

import (
	platform "github.com/arangodb-managed/apis/platform/v1"
)

// Region returns a single region formatted for humans.
func Region(x *platform.Region, opts Options) string {
	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"provider-id", x.GetProviderId()},
		kv{"location", x.GetLocation()},
		kv{"available", formatBool(opts, x.GetAvailable())},
	)
}

// RegionList returns a list of regions formatted for humans.
func RegionList(list []*platform.Region, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			kv{"id", x.GetId()},
			kv{"provider-id", x.GetProviderId()},
			kv{"location", x.GetLocation()},
			kv{"available", formatBool(opts, x.GetAvailable())},
		}
	}, false)
}
