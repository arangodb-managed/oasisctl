//
// DISCLAIMER
//
// Copyright 2019 ArangoDB GmbH, Cologne, Germany
//
// Author Ewout Prangsma
//

package format

import (
	iam "github.com/arangodb-managed/apis/iam/v1"
)

// User returns a single user formatted for humans.
func User(x *iam.User, opts Options) string {
	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"name", x.GetName()},
		kv{"email", x.GetEmail()},
		kv{"created_at", formatTime(opts, x.GetCreatedAt(), "-")},
		kv{"last_login_at", formatTime(opts, x.GetLastLoginAt(), "-")},
		kv{"last_ip", x.GetLastIp()},
	)
}

// UserList returns a list of users formatted for humans.
func UserList(list []*iam.User, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			kv{"id", x.GetId()},
			kv{"name", x.GetName()},
			kv{"email", x.GetEmail()},
			kv{"created_at", formatTime(opts, x.GetCreatedAt(), "-")},
			kv{"last_login_at", formatTime(opts, x.GetLastLoginAt(), "-")},
			kv{"last_ip", x.GetLastIp()},
		}
	}, false)
}
