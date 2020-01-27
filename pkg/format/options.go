//
// DISCLAIMER
//
// Copyright 2019 ArangoDB GmbH, Cologne, Germany
//
// Author Ewout Prangsma
//

package format

const (
	// DefaultFormat specifies default value for Options.Format
	DefaultFormat = formatTable

	formatTable = "table"
	formatJSON  = "json"
)

// Options that control the formatter.
type Options struct {
	Format string
}
