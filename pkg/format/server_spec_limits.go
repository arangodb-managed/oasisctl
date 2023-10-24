//
// DISCLAIMER
//
// Copyright 2020-2023 ArangoDB GmbH, Cologne, Germany
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

package format

import (
	"fmt"
	"strings"

	data "github.com/arangodb-managed/apis/data/v1"
)

// ServersSpecLimits returns a single server specification limts formatted for humans.
func ServersSpecLimits(x *data.ServersSpecLimits, opts Options) string {
	return formatObject(opts,
		kv{"coordinators", serversSpecLimitsLimits(x.GetCoordinators(), "")},
		kv{"coordinator-memory-size", serversSpecLimitsLimits(x.GetCoordinatorMemorySize(), "GiB")},
		kv{"dbservers", serversSpecLimitsLimits(x.GetDbservers(), "")},
		kv{"dbserver-memory-size", serversSpecLimitsLimits(x.GetDbserverMemorySize(), "GiB")},
		kv{"dbserver-disk-size", serversSpecLimitsLimits(x.GetDbserverDiskSize(), "GiB")},
	)
}

// serversSpecLimitsLimits returns a single property  specification limts formatted for humans.
func serversSpecLimitsLimits(x *data.ServersSpecLimits_Limits, unit string) string {
	if av := x.GetAllowedValues(); len(av) > 0 {
		list := make([]string, 0, len(av))
		for _, v := range av {
			list = append(list, fmt.Sprintf("%d%s", v, unit))
		}
		return strings.Join(list, ", ")
	}
	return fmt.Sprintf("%d%s - %d%s", x.GetMin(), unit, x.GetMax(), unit)
}
