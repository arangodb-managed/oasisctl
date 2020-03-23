//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
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
