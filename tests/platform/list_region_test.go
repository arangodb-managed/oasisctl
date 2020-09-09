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

package platform

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arangodb-managed/oasisctl/tests"
)

func TestListRegion(t *testing.T) {
	args := []string{"list", "regions", "--provider-id=aks"}
	compare := `Id                | Provider-Id | Location                 | Available
aks-canadacentral | aks         | Central Canada, Toronto  | ✓
aks-eastus2       | aks         | East US, Virginia        | ✓
aks-uksouth       | aks         | UK, London               | ✓
aks-westeurope    | aks         | West Europe, Netherlands | ✓
aks-westus2       | aks         | West US, Washington      | ✓
$`
	out, err := tests.RunCommand(args)
	require.NoError(t, err)
	assert.True(t, tests.CompareOutput(out, []byte(compare)))
}
