//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package format

// formatBool returns a human readable checkmark for the given boolean
func formatBool(x bool) string {
	if x {
		return "x"
	}
	return "-"
}
