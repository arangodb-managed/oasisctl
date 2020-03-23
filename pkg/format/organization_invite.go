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

// OrganizationInvite returns a single organization member formatted for humans.
func OrganizationInvite(ctx context.Context, x *rm.OrganizationInvite, iamc iam.IAMServiceClient, opts Options) string {
	userName := "-"
	if x.GetUserId() != "" {
		if user, err := iamc.GetUser(ctx, &common.IDOptions{Id: x.GetUserId()}); err == nil {
			userName = user.GetName()
		}
	}
	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"email", x.GetEmail()},
		kv{"organization", x.GetOrganizationName()},
		kv{"created-by", x.GetCreatedByName()},
		kv{"accepted", formatTime(opts, x.GetAcceptedAt(), "-")},
		kv{"rejected", formatTime(opts, x.GetRejectedAt(), "-")},
		kv{"user", userName},
		kv{"created-at", formatTime(opts, x.GetCreatedAt())},
		kv{"url", x.GetUrl()},
	)
}

// OrganizationInviteList returns a list of organization members formatted for humans.
func OrganizationInviteList(ctx context.Context, list []*rm.OrganizationInvite, iamc iam.IAMServiceClient, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		userName := "-"
		if x.GetUserId() != "" {
			if user, err := iamc.GetUser(ctx, &common.IDOptions{Id: x.GetUserId()}); err == nil {
				userName = user.GetName()
			}
		}
		return []kv{
			kv{"id", x.GetId()},
			kv{"email", x.GetEmail()},
			kv{"organization", x.GetOrganizationName()},
			kv{"created-by", x.GetCreatedByName()},
			kv{"accepted", formatTime(opts, x.GetAcceptedAt(), "-")},
			kv{"rejected", formatTime(opts, x.GetRejectedAt(), "-")},
			kv{"user", userName},
			kv{"created-at", formatTime(opts, x.GetCreatedAt())},
			kv{"url", x.GetUrl()},
		}
	}, false)
}
