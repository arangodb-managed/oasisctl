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

// Organization returns a single organization formatted for humans.
func Organization(x *rm.Organization, opts Options) string {
	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"name", x.GetName()},
		kv{"description", x.GetDescription()},
		kv{"url", x.GetUrl()},
		kv{"created-at", formatTime(x.GetCreatedAt())},
		kv{"deleted-at", formatTime(x.GetDeletedAt(), "-")},
	)
}

// OrganizationList returns a list of organizations formatted for humans.
func OrganizationList(list []*rm.Organization, opts Options) string {
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
