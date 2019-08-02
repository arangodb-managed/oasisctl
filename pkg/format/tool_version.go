//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package format

// CLIVersion returns a single version formatted for humans.
func CLIVersion(version string, opts Options) string {
	return formatObject(opts,
		kv{"version", version},
	)
}
