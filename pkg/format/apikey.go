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

// APIKey returns a single api key formatted for humans.
func APIKey(x *iam.APIKey, opts Options) string {
	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"user-id", x.GetUserId()},
		kv{"organization-id", x.GetOrganizationId()},
		kv{"readonly", formatBool(opts, x.GetIsReadonly())},
		kv{"created-at", formatTime(opts, x.GetCreatedAt())},
		kv{"expires-at", formatTime(opts, x.GetExpiresAt())},
		kv{"revoked-at", formatTime(opts, x.GetRevokedAt())},
	)
}

// APIKeyList returns a list of api keys formatted for humans.
func APIKeyList(list []*iam.APIKey, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			kv{"id", x.GetId()},
			kv{"user-id", x.GetUserId()},
			kv{"organization-id", x.GetOrganizationId()},
			kv{"readonly", formatBool(opts, x.GetIsReadonly())},
			kv{"created-at", formatTime(opts, x.GetCreatedAt())},
			kv{"expires-at", formatTime(opts, x.GetExpiresAt())},
			kv{"revoked-at", formatTime(opts, x.GetRevokedAt())},
		}
	}, false)
}
