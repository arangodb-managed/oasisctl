//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package format

import (
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

// Project returns a single project formatted for humans.
func Project(x *rm.Project, opts Options) string {
	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"name", x.GetName()},
		kv{"description", x.GetDescription()},
		kv{"url", x.GetUrl()},
		kv{"created-at", formatTime(opts, x.GetCreatedAt())},
		kv{"deleted-at", formatTime(opts, x.GetDeletedAt(), "-")},
	)
}

// ProjectList returns a list of projects formatted for humans.
func ProjectList(list []*rm.Project, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			kv{"id", x.GetId()},
			kv{"name", x.GetName()},
			kv{"description", x.GetDescription()},
			kv{"url", x.GetUrl()},
			kv{"created-at", formatTime(opts, x.GetCreatedAt())},
		}
	}, false)
}
