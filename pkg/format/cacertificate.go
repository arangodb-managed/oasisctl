//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package format

import (
	crypto "github.com/arangodb-managed/apis/crypto/v1"
)

// CACertificate returns a single ca certificate formatted for humans.
func CACertificate(x *crypto.CACertificate, opts Options) string {
	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"name", x.GetName()},
		kv{"description", x.GetDescription()},
		kv{"lifetime", formatDuration(x.GetLifetime())},
		kv{"url", x.GetUrl()},
		kv{"created-at", formatTime(x.GetCreatedAt())},
		kv{"deleted-at", formatTime(x.GetDeletedAt(), "-")},
	)
}

// CACertificateList returns a list of ca certificates formatted for humans.
func CACertificateList(list []*crypto.CACertificate, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			kv{"id", x.GetId()},
			kv{"name", x.GetName()},
			kv{"description", x.GetDescription()},
			kv{"lifetime", formatDuration(x.GetLifetime())},
			kv{"url", x.GetUrl()},
			kv{"created-at", formatTime(x.GetCreatedAt())},
		}
	})
}
