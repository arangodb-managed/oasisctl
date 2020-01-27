//
// DISCLAIMER
//
// Copyright 2019 ArangoDB GmbH, Cologne, Germany
//
// Author Ewout Prangsma
//

package format

import (
	"strconv"
)

// formatBool returns a human readable checkmark for the given boolean
func formatBool(opts Options, x bool) string {
	if opts.Format == formatJSON {
		return strconv.FormatBool(x)
	}
	if x {
		return "x"
	}
	return "-"
}
