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

func TestCreateCrypto(t *testing.T) {
	org, err := tests.GetDefaultOrganization()
	require.NoError(t, err)
	proj, err := tests.GetDefaultProject(org)
	require.NoError(t, err)
	certName := "testcertificate"
	args := []string{"create", "cacertificate", "--name=" + certName, "--organization-id=" + org, "--project-id=" + proj}

	compare := `^Success!
Id                         .*
Name                       testcertificate
Description                
Lifetime                   \d+h0m0s
Url                        /Organization/\d+/Project/\d+/CACertificate/.*
Use-Well-Known-Certificate -
Created-At                 now
Deleted-At                 -
$`
	out, err := tests.RunCommand(args)
	require.NoError(t, err)

	certId, err := tests.GetResourceID(string(out))
	require.NoError(t, err)
	defer func() {
		conn := cmd.MustDialAPI()
		defer conn.Close()
		ctx := cmd.ContextWithToken()
		cryptoc := crypto.NewCryptoServiceClient(conn)
		if _, err := cryptoc.DeleteCACertificate(ctx, &common.IDOptions{Id: certId}); err != nil {
			t.Log(err)
		}
	}()

	assert.True(t, tests.CompareOutput(out, []byte(compare)))
}
