//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package format

import (
	iam "github.com/arangodb-managed/apis/iam/v1"
)

// Group returns a single group formatted for humans.
func Group(x *iam.Group, opts Options) string {
	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"name", x.GetName()},
		kv{"description", x.GetDescription()},
		kv{"url", x.GetUrl()},
		kv{"created-at", formatTime(x.GetCreatedAt())},
		kv{"deleted-at", formatTime(x.GetDeletedAt(), "-")},
	)
}

// GroupList returns a list of groups formatted for humans.
func GroupList(list []*iam.Group, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			kv{"id", x.GetId()},
			kv{"name", x.GetName()},
			kv{"description", x.GetDescription()},
			kv{"url", x.GetUrl()},
			kv{"created-at", formatTime(x.GetCreatedAt())},
		}
	})
}
