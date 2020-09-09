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
// +build e2e

package resourcemanager

import (
	"testing"

	"github.com/arangodb-managed/oasisctl/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"
	"github.com/arangodb-managed/oasisctl/cmd"
)

func TestCreateApiKey(t *testing.T) {
	cmd.RootCmd.PersistentPreRun(nil, nil)
	ctx := cmd.ContextWithToken()
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	org, err := tests.GetDefaultOrganization()
	require.NoError(t, err)

	defer func() {
		list, err := iamc.ListAPIKeys(ctx, &common.ListOptions{ContextId: org})
		if err != nil {
			t.Log(err)
		}

		// We remove all keys.
		for _, key := range list.GetItems() {
			if _, err := iamc.DeleteAPIKey(ctx, &common.IDOptions{Id: key.GetId()}); err != nil {
				t.Log(err)
			}
		}
	}()

	compare := `^Success!
Success!
Id     .*
Secret .*
$`

	args := []string{"create", "apikey"}
	out, err := tests.RunCommand(args)
	require.NoError(t, err)
	assert.True(t, tests.CompareOutput(out, []byte(compare)))
}
