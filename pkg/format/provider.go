//
// DISCLAIMER
//
// Copyright 2019 ArangoDB GmbH, Cologne, Germany
//
// Author Ewout Prangsma
//

package format

import (
	platform "github.com/arangodb-managed/apis/platform/v1"
)

// Provider returns a single provider formatted for humans.
func Provider(x *platform.Provider, opts Options) string {
	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"name", x.GetName()},
	)
}

// ProviderList returns a list of providers formatted for humans.
func ProviderList(list []*platform.Provider, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			kv{"id", x.GetId()},
			kv{"name", x.GetName()},
		}
	}, false)
}
