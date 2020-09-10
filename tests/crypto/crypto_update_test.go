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

package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	common "github.com/arangodb-managed/apis/common/v1"
	crypto "github.com/arangodb-managed/apis/crypto/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/tests"
)

func TestUpdateCrypto(t *testing.T) {
	cmd.RootCmd.PersistentPreRun(nil, nil)
	ctx := cmd.ContextWithToken()
	cryptoc, project, err := getCryptoClientAndProject()
	require.NoError(t, err)
	org, err := tests.GetDefaultOrganization()
	require.NoError(t, err)
	// Create a certificate via the api.
	result, err := cryptoc.CreateCACertificate(ctx, &crypto.CACertificate{
		ProjectId: project,
		Name:      "TestUpdateCrypto",
	})
	require.NoError(t, err)

	// Cleanup
	defer func() {
		if _, err := cryptoc.DeleteCACertificate(ctx, &common.IDOptions{Id: result.GetId()}); err != nil {
			t.Log("Failed to cleanup certificate: ", err)
		}
	}()

	args := []string{"update", "cacertificate", "--cacertificate-id=" + result.GetId(), "--name=NewName", "--description=NewDescription", "--use-well-known-certificate=true", "--organization-id=" + org, "--project-id=" + project}
	compare := `Updated CA certificate!
Id                         ` + result.GetId() + `
Name                       NewName
Description                NewDescription
Lifetime                   \d+h0m0s
Url                        ` + result.GetUrl() + `
Use-Well-Known-Certificate ✓
Created-At                 .*
Deleted-At                 -
$`
	out, err := tests.RunCommand(args)
	require.NoError(t, err)
	assert.True(t, tests.CompareOutput(out, []byte(compare)))
}
