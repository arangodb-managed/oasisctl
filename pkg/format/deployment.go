//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package format

import (
	"fmt"

	data "github.com/arangodb-managed/apis/data/v1"
)

// Deployment returns a single deployment formatted for humans.
func Deployment(x *data.Deployment, opts Options, showRootpassword bool) string {
	pwd := func(p string) string {
		if showRootpassword {
			return p
		}
		return "*** use '--show-root-password' to expose ***"
	}

	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"name", x.GetName()},
		kv{"description", x.GetDescription()},
		kv{"region", x.GetRegionId()},
		kv{"version", x.GetVersion()},
		kv{"ipwhitelist", x.GetIpwhitelistId()},
		kv{"url", x.GetUrl()},
		kv{"created-at", formatTime(opts, x.GetCreatedAt())},
		kv{"deleted-at", formatTime(opts, x.GetDeletedAt(), "-")},
		kv{"expires-at", formatTime(opts, x.GetExpiration().GetExpiresAt(), "-")},

		kv{"coordinators", x.GetServers().GetCoordinators()},
		kv{"coordinator-memory-size", fmt.Sprintf("%d%s", x.GetServers().GetCoordinatorMemorySize(), "GB")},
		kv{"dbservers", x.GetServers().GetDbservers()},
		kv{"dbserver-memory-size", fmt.Sprintf("%d%s", x.GetServers().GetDbserverMemorySize(), "GB")},
		kv{"dbserver-disk-size", fmt.Sprintf("%d%s", x.GetServers().GetDbserverDiskSize(), "GB")},

		kv{"bootstrapped-at", formatTime(opts, x.GetStatus().GetBootstrappedAt(), "-")},
		kv{"endpoint-url", x.GetStatus().GetEndpoint()},
		kv{"root-password", pwd(x.GetAuthentication().GetRootPassword())},
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
			kv{"version", x.GetVersion()},
			kv{"ipwhitelist", x.GetIpwhitelistId()},
			kv{"url", x.GetUrl()},
			kv{"created-at", formatTime(opts, x.GetCreatedAt())},
		}
	}, false)
}
