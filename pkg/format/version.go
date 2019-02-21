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

// Version returns a single version formatted for humans.
func Version(x *data.Version, opts Options) string {
	return formatObject(opts,
		kv{"version", x.GetVersion()},
	)
}

// VersionList returns a list of versions formatted for humans.
func VersionList(list []*data.Version, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			kv{"version", x.GetVersion()},
		}
	})
}
