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

	crypto "github.com/arangodb-managed/apis/crypto/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/tests"
)

func TestDeleteCertificate(t *testing.T) {
	// Initialize the root command.
	cmd.RootCmd.PersistentPreRun(nil, nil)
	ctx := cmd.ContextWithToken()
	conn := cmd.MustDialAPI()
	defer conn.Close()
	cryptoc := crypto.NewCryptoServiceClient(conn)
	org, err := tests.GetDefaultOrganization()
	require.NoError(t, err)
	project, err := tests.GetDefaultProject(org)
	require.NoError(t, err)
	// Create a certificate via the api.
	result, err := cryptoc.CreateCACertificate(ctx, &crypto.CACertificate{
		ProjectId: project,
		Name:      "TestDeleteCertificate",
	})
	require.NoError(t, err)
	args := []string{"delete", "cacertificate", "--cacertificate-id=" + result.GetId(), "--organization-id=" + org, "--project-id=" + project}
	compare := `^Deleted.CA.certificate!
$`
	out, err := tests.RunCommand(args)
	require.NoError(t, err)
	assert.True(t, tests.CompareOutput(out, []byte(compare)))

	// Try getting the deleted certificate
	args = []string{"get", "cacertificate", "--cacertificate-id=" + result.GetId(), "--organization-id=" + org, "--project-id=" + project}
	_, err = tests.RunCommand(args)
	require.Error(t, err)
}
