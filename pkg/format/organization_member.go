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

// OrganizationMember returns a single organization member formatted for humans.
func OrganizationMember(ctx context.Context, x *rm.Member, iamc iam.IAMServiceClient, opts Options) string {
	userName := "?"
	userEmail := "?"
	user, err := iamc.GetUser(ctx, &common.IDOptions{Id: x.GetUserId()})
	if err == nil {
		userName = user.GetName()
		userEmail = user.GetEmail()
	}
	return formatObject(opts,
		kv{"id", x.GetUserId()},
		kv{"name", userName},
		kv{"email", userEmail},
		kv{"owner", formatBool(opts, x.GetOwner())},
	)
}

// OrganizationMemberList returns a list of organization members formatted for humans.
func OrganizationMemberList(ctx context.Context, list []*rm.Member, iamc iam.IAMServiceClient, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		userName := "?"
		userEmail := "?"
		user, err := iamc.GetUser(ctx, &common.IDOptions{Id: x.GetUserId()})
		if err == nil {
			userName = user.GetName()
			userEmail = user.GetEmail()
		}
		return []kv{
			kv{"id", x.GetUserId()},
			kv{"name", userName},
			kv{"email", userEmail},
			kv{"owner", formatBool(opts, x.GetOwner())},
		}
	}, false)
}
