//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package format

import (
	"context"
	"sort"
	"strings"

	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"
)

// Policy returns a single policy formatted for humans.
func Policy(ctx context.Context, x *iam.Policy, iamc iam.IAMServiceClient, opts Options) string {
	list := x.GetBindings()
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		roleName := x.GetRoleId()
		var permissions []string
		if role, err := iamc.GetRole(ctx, &common.IDOptions{Id: x.GetRoleId()}); err == nil {
			roleName = role.GetName()
			permissions = role.GetPermissions()
			sort.Strings(permissions)
		}
		return []kv{
			kv{"id", x.GetId()},
			kv{"member-id", x.GetMemberId()},
			kv{"role", roleName},
			kv{"delete-not-allowed", formatBool(opts, x.GetDeleteNotAllowed())},
			kv{"permissions", strings.Join(permissions, ", ")},
		}
	}, false)
}
