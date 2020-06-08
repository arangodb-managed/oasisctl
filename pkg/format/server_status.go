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
	data "github.com/arangodb-managed/apis/data/v1"
)

// DeploymentList returns a list of deployments formatted for humans.
func ServerStatusList(list []*data.Deployment_ServerStatus, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		d := []kv{
			{"id", x.GetId()},
			{"description", x.GetDescription()},
			{"version", x.GetVersion()},
			{"creating", x.GetType()},
			{"created-at", formatTime(opts, x.GetCreatedAt())},
			{"last-started-at", formatTime(opts, x.GetLastStartedAt())},
			{"creating", x.GetCreating()},
			{"ready", x.GetReady()},
			{"failed", x.GetFailed()},
			{"upgrading", x.GetUpgrading()},
			{"ok", x.GetOk()},
			{"member-of-cluster", x.GetMemberOfCluster()},
		}
		return d
	}, false)
}
