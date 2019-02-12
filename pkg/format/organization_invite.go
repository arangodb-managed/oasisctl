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

	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

// OrganizationInvite returns a single organization member formatted for humans.
func OrganizationInvite(ctx context.Context, x *rm.OrganizationInvite, iamc iam.IAMServiceClient, opts Options) string {
	createdByUserName := "?"
	if user, err := iamc.GetUser(ctx, &common.IDOptions{Id: x.GetCreatedById()}); err == nil {
		createdByUserName = user.GetName()
	}
	userName := "-"
	if x.GetUserId() != "" {
		if user, err := iamc.GetUser(ctx, &common.IDOptions{Id: x.GetUserId()}); err == nil {
			userName = user.GetName()
		}
	}
	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"email", x.GetEmail()},
		kv{"created-by", createdByUserName},
		kv{"accepted", formatTime(x.GetAcceptedAt(), "-")},
		kv{"rejected", formatTime(x.GetRejectedAt(), "-")},
		kv{"user", userName},
		kv{"created-at", formatTime(x.GetCreatedAt())},
	)
}

// OrganizationInviteList returns a list of organization members formatted for humans.
func OrganizationInviteList(ctx context.Context, list []*rm.OrganizationInvite, iamc iam.IAMServiceClient, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		createdByUserName := "?"
		if user, err := iamc.GetUser(ctx, &common.IDOptions{Id: x.GetCreatedById()}); err == nil {
			createdByUserName = user.GetName()
		}
		userName := "-"
		if x.GetUserId() != "" {
			if user, err := iamc.GetUser(ctx, &common.IDOptions{Id: x.GetUserId()}); err == nil {
				userName = user.GetName()
			}
		}
		return []kv{
			kv{"id", x.GetId()},
			kv{"email", x.GetEmail()},
			kv{"created-by", createdByUserName},
			kv{"accepted", formatTime(x.GetAcceptedAt(), "-")},
			kv{"rejected", formatTime(x.GetRejectedAt(), "-")},
			kv{"user", userName},
			kv{"created-at", formatTime(x.GetCreatedAt())},
		}
	})
}
