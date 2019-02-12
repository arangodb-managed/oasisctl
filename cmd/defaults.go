//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package cmd

func defaultOrganization() string       { return envOrDefault("ORGANIZATION", "") }
func defaultOrganizationInvite() string { return envOrDefault("ORGANIZATION_INVITE", "") }
func defaultProject() string            { return envOrDefault("PROJECT", "") }
func defaultGroup() string              { return envOrDefault("GROUP", "") }
func defaultRole() string               { return envOrDefault("ROLE", "") }
func defaultURL() string                { return envOrDefault("URL", "") }
