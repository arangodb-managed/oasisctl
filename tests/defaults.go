//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Author Gergely Brautigam
//

package tests

import (
	"errors"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
)

// GetDefaultOrganization returns the first organization for the user to work with.
func GetDefaultOrganization() (string, error) {
	cmd.RootCmd.PersistentPreRun(nil, nil)
	org := cmd.DefaultOrganization()
	if org != "" {
		return org, nil
	}
	ctx := cmd.ContextWithToken()
	conn := cmd.MustDialAPI()

	rmc := rm.NewResourceManagerServiceClient(conn)

	// Get the first organization if default is not set.
	list, err := rmc.ListOrganizations(ctx, &common.ListOptions{})
	if err != nil {
		return "", err
	}
	if len(list.Items) == 0 {
		return "", errors.New("no organizations found")
	}
	return list.GetItems()[0].GetId(), nil
}

// GetDefaultProject returns the first project for the user to work with.
func GetDefaultProject(org string) (string, error) {
	cmd.RootCmd.PersistentPreRun(nil, nil)
	proj := cmd.DefaultProject()
	if proj != "" {
		return proj, nil
	}
	ctx := cmd.ContextWithToken()
	conn := cmd.MustDialAPI()

	rmc := rm.NewResourceManagerServiceClient(conn)

	// Get the first project if default is not set.
	list, err := rmc.ListProjects(ctx, &common.ListOptions{ContextId: org})
	if err != nil {
		return "", err
	}
	if len(list.Items) == 0 {
		return "", errors.New("no projects found")
	}
	return list.GetItems()[0].GetId(), nil
}
