//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package format

import (
	"strings"

	security "github.com/arangodb-managed/apis/security/v1"
)

// IPWhitelist returns a single IP whitelist formatted for humans.
func IPWhitelist(x *security.IPWhitelist, opts Options) string {
	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"name", x.GetName()},
		kv{"description", x.GetDescription()},
		kv{"cidr-ranges", strings.Join(x.GetCidrRanges(), ", ")},
		kv{"url", x.GetUrl()},
		kv{"created-at", formatTime(x.GetCreatedAt())},
	)
}

// IPWhitelistList returns a list of IP whitelists formatted for humans.
func IPWhitelistList(list []*security.IPWhitelist, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			kv{"id", x.GetId()},
			kv{"name", x.GetName()},
			kv{"description", x.GetDescription()},
			kv{"cidr-ranges", strings.Join(x.GetCidrRanges(), ", ")},
			kv{"url", x.GetUrl()},
			kv{"created-at", formatTime(x.GetCreatedAt())},
		}
	}, false)
}
