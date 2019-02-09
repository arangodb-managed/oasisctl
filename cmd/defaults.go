//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package cmd

func defaultOrganization() string { return envOrDefault("ORGANIZATION", "") }
func defaultProject() string      { return envOrDefault("PROJECT", "") }
