//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

// +build !windows

package cmd

func supportsColor() bool {
	return true
}
