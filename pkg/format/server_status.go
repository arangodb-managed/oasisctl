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
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"

	data "github.com/arangodb-managed/apis/data/v1"
)

func formatCPU(value float64) string {
	if value < 1.0 {
		return strconv.Itoa(int(value*1000.0)) + "m"
	}
	return humanize.FormatFloat("", value)
}

func formatDisk(value uint64, serverType string) string {
	if value == 0 && serverType == "Coordinator" {
		return "-"
	}
	return humanize.Bytes(value)
}

// ServerStatusList returns a list of deployment servers formatted for humans.
func ServerStatusList(list []*data.Deployment_ServerStatus, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		d := []kv{
			{"id", x.GetId()},
			{"description", x.GetDescription()},
			{"version", x.GetVersion()},
			{"type", x.GetType()},
			{"memory", humanize.Bytes(uint64(x.GetLastMemoryUsage()))},
			{"cpu", formatCPU(float64(x.GetLastCpuUsage()))},
			{"disk", formatDisk(uint64(x.GetDataVolumeInfo().GetUsedBytes()), x.GetType())},
			{"created-at", formatTime(opts, x.GetCreatedAt())},
			{"last-started-at", formatTime(opts, x.GetLastStartedAt())},
			{"creating", formatBool(opts, x.GetCreating())},
			{"ready", formatBool(opts, x.GetReady())},
			{"failed", formatBool(opts, x.GetFailed())},
			{"upgrading", formatBool(opts, x.GetUpgrading())},
			{"ok", formatBool(opts, x.GetOk())},
			{"member-of-cluster", formatBool(opts, x.GetMemberOfCluster())},
		}
		return d
	}, false)
}

// ServerStatusListAsRows returns a list of deployment servers formatted for a table.
func ServerStatusListAsRows(list []*data.Deployment_ServerStatus, opts Options) []string {
	all := ServerStatusList(list, opts)
	return strings.Split(all, "\n")
}
