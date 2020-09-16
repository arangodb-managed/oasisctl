//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
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
	defer conn.Close()

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
	defer conn.Close()

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
