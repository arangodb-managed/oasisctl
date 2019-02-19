//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

// +build !windows

package format

// formatBool returns a human readable checkmark for the given boolean
func formatBool(x bool) string {
	if x {
		return "\u2713"
	}
	return "-"
}
