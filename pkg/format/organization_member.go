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

	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

// OrganizationMember returns a single organization member formatted for humans.
func OrganizationMember(ctx context.Context, x *rm.Member, iamc iam.IAMServiceClient, opts Options) string {
	userName := "?"
	userEmail := "?"
	userCreatedAt := "?"
	userLastLoginAt := "?"
	userLastIP := "?"
	user, err := iamc.GetUser(ctx, &common.IDOptions{Id: x.GetUserId()})
	if err == nil {
		userName = user.GetName()
		userEmail = user.GetEmail()
		userCreatedAt = formatTime(opts, user.GetCreatedAt(), "?")
		userLastLoginAt = formatTime(opts, user.GetLastLoginAt(), "-")
		userLastIP = user.GetLastIp()
	}
	return formatObject(opts,
		kv{"id", x.GetUserId()},
		kv{"name", userName},
		kv{"email", userEmail},
		kv{"created_at", userCreatedAt},
		kv{"last_login_at", userLastLoginAt},
		kv{"last_ip", userLastIP},
		kv{"owner", formatBool(opts, x.GetOwner())},
	)
}

// OrganizationMemberList returns a list of organization members formatted for humans.
func OrganizationMemberList(ctx context.Context, list []*rm.Member, iamc iam.IAMServiceClient, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		userName := "?"
		userEmail := "?"
		userCreatedAt := "?"
		userLastLoginAt := "?"
		userLastIP := "?"
		user, err := iamc.GetUser(ctx, &common.IDOptions{Id: x.GetUserId()})
		if err == nil {
			userName = user.GetName()
			userEmail = user.GetEmail()
			userCreatedAt = formatTime(opts, user.GetCreatedAt(), "?")
			userLastLoginAt = formatTime(opts, user.GetLastLoginAt(), "-")
			userLastIP = user.GetLastIp()
		}
		return []kv{
			kv{"id", x.GetUserId()},
			kv{"name", userName},
			kv{"email", userEmail},
			kv{"created_at", userCreatedAt},
			kv{"last_login_at", userLastLoginAt},
			kv{"last_ip", userLastIP},
			kv{"owner", formatBool(opts, x.GetOwner())},
		}
	}, false)
}
