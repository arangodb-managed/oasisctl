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
)

// GroupMember returns a single organization member formatted for humans.
func GroupMember(ctx context.Context, x string, iamc iam.IAMServiceClient, opts Options) string {
	userName := "?"
	userEmail := "?"
	user, err := iamc.GetUser(ctx, &common.IDOptions{Id: x})
	if err == nil {
		userName = user.GetName()
		userEmail = user.GetEmail()
	}
	return formatObject(opts,
		kv{"id", x},
		kv{"name", userName},
		kv{"email", userEmail},
	)
}

// GroupMemberList returns a list of group members formatted for humans.
func GroupMemberList(ctx context.Context, list []string, iamc iam.IAMServiceClient, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		userName := "?"
		userEmail := "?"
		user, err := iamc.GetUser(ctx, &common.IDOptions{Id: x})
		if err == nil {
			userName = user.GetName()
			userEmail = user.GetEmail()
		}
		return []kv{
			kv{"id", x},
			kv{"name", userName},
			kv{"email", userEmail},
		}
	})
}
