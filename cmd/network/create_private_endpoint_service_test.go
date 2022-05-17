//
// DISCLAIMER
//
// Copyright 2021-2022 ArangoDB GmbH, Cologne, Germany
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

package network

import (
	"testing"

	"github.com/stretchr/testify/assert"

	network "github.com/arangodb-managed/apis/network/v1"
)

func TestParseAwsPrincipal(t *testing.T) {
	result, err := parseAwsPrincipal("AccountId")
	assert.NoError(t, err)
	assert.Equal(t, "AccountId", result.AccountID)
	assert.Empty(t, result.RoleName)
	assert.Empty(t, result.UserName)

	result, err = parseAwsPrincipal("AccountId/Role/TestRole")
	assert.NoError(t, err)
	assert.Equal(t, "AccountId", result.AccountID)
	assert.Equal(t, "TestRole", result.RoleName)
	assert.Empty(t, result.UserName)

	result, err = parseAwsPrincipal("AccountId/User/TestUser")
	assert.NoError(t, err)
	assert.Equal(t, "AccountId", result.AccountID)
	assert.Empty(t, result.RoleName)
	assert.Equal(t, "TestUser", result.UserName)

	_, err = parseAwsPrincipal("AccountId/Error1")
	assert.Error(t, err)

	_, err = parseAwsPrincipal("AccountId/Error/Something")
	assert.Error(t, err)
}

func TestGetAwsPrincipals(t *testing.T) {
	result, err := getAwsPrincipals([]string{
		"a1/Role/r1",
		"a1/Role/r2",
		"a1/Role/r3",
		"a1/User/u1",
		"a1/User/u2",
		"a1/User/u3",
		"a2",
		"a3/Role/r1",
	})
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.ElementsMatch(t, result, []*network.PrivateEndpointService_AwsPrincipals{
		{
			AccountId: "a1",
			RoleNames: []string{"r1", "r2", "r3"},
			UserNames: []string{"u1", "u2", "u3"},
		},
		{
			AccountId: "a2",
		},
		{
			AccountId: "a3",
			RoleNames: []string{"r1"},
		},
	})
}
