//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package format

// GroupMember returns a single organization member formatted for humans.
func GroupMember(x string, opts Options) string {
	return formatObject(opts,
		kv{"user-id", x},
	)
}

// GroupMemberList returns a list of group members formatted for humans.
func GroupMemberList(list []string, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			kv{"user-id", x},
		}
	})
}
