//
// DISCLAIMER
//
// Copyright 2019 ArangoDB GmbH, Cologne, Germany
//
// Author Ewout Prangsma
//

// +build !windows

package cmd

func supportsColor() bool {
	return true
}
