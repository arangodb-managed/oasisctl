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

// OrganizationMember returns a single organization member formatted for humans.
func OrganizationMember(x *rm.Member, opts Options) string {
	return formatObject(opts,
		kv{"user-id", x.GetUserId()},
		kv{"owner", formatBool(x.GetOwner())},
	)
}

// OrganizationMemberList returns a list of organization members formatted for humans.
func OrganizationMemberList(list []*rm.Member, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			kv{"user-id", x.GetUserId()},
			kv{"owner", formatBool(x.GetOwner())},
		}
	})
}
