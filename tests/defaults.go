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

	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	common "github.com/arangodb-managed/apis/common/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
)

// GetDefaultOrganization returns the first organization for the user to work with.
func GetDefaultOrganization() (string, error) {
	cmd.RootCmd.PersistentPreRun(nil, nil)
	ctx := cmd.ContextWithToken()
	conn := cmd.MustDialAPI()
	org := cmd.DefaultOrganization()
	if org != "" {
		return org, nil
	}

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
