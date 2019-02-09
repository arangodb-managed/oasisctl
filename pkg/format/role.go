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

	iam "github.com/arangodb-managed/apis/iam/v1"
)

// Role returns a single role formatted for humans.
func Role(x *iam.Role, opts Options) string {
	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"name", x.GetName()},
		kv{"description", x.GetDescription()},
		kv{"predefined", x.GetIsPredefined()},
		kv{"permissions", strings.Join(x.GetPermissions(), ", ")},
		//kv{"url", x.GetUrl()},
		kv{"created-at", formatTime(x.GetCreatedAt())},
		kv{"deleted-at", formatTime(x.GetCreatedAt(), "-")},
	)
}

// RoleList returns a list of roles formatted for humans.
func RoleList(list []*iam.Role, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			kv{"id", x.GetId()},
			kv{"name", x.GetName()},
			kv{"description", x.GetDescription()},
			kv{"predefined", x.GetIsPredefined()},
			kv{"permissions", strings.Join(x.GetPermissions(), ", ")},
			//kv{"url", x.GetUrl()},
			kv{"created-at", formatTime(x.GetCreatedAt())},
		}
	})
}
