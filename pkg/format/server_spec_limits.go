//
// DISCLAIMER
//
// Copyright 2019 ArangoDB GmbH, Cologne, Germany
//
// Author Ewout Prangsma
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
		kv{"coordinator-memory-size", serversSpecLimitsLimits(x.GetCoordinatorMemorySize(), "GB")},
		kv{"dbservers", serversSpecLimitsLimits(x.GetDbservers(), "")},
		kv{"dbserver-memory-size", serversSpecLimitsLimits(x.GetDbserverMemorySize(), "GB")},
		kv{"dbserver-disk-size", serversSpecLimitsLimits(x.GetDbserverDiskSize(), "GB")},
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
