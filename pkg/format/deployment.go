//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package format

import (
	data "github.com/arangodb-managed/apis/data/v1"
)

// Deployment returns a single deployment formatted for humans.
func Deployment(x *data.Deployment, opts Options) string {
	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"name", x.GetName()},
		kv{"description", x.GetDescription()},
		kv{"region", x.GetRegionId()},
		kv{"url", x.GetUrl()},
		kv{"created-at", formatTime(x.GetCreatedAt())},
		kv{"deleted-at", formatTime(x.GetDeletedAt(), "-")},
		// TODO other relevant fields
	)
}

// DeploymentList returns a list of deployments formatted for humans.
func DeploymentList(list []*data.Deployment, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			kv{"id", x.GetId()},
			kv{"name", x.GetName()},
			kv{"description", x.GetDescription()},
			kv{"region", x.GetRegionId()},
			kv{"url", x.GetUrl()},
			kv{"created-at", formatTime(x.GetCreatedAt())},
		}
	}, false)
}
