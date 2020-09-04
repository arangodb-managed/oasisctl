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

package crypto

import (
	"testing"

	_ "github.com/arangodb-managed/oasisctl/cmd/crypto"
	"github.com/arangodb-managed/oasisctl/tests"
)

func TestCreateCrypto(t *testing.T) {
	args := []string{"create", "cacertificate", "--name=testcertificate"}
	compare := `^Success!
Id                         .*
Name                       testcertificate
Description                
Lifetime                   \d+h0m0s
Url                        /Organization/_support/Project/\d+/CACertificate/.*
Use-Well-Known-Certificate -
Created-At                 now
Deleted-At                 -
$`
	// perform any setups in here.
	out, err := tests.RunCommand(args)
	if err != nil {
		t.Fatal(err)
	}
	if !tests.CompareOutput(out, []byte(compare)) {
		t.FailNow()
	}
}
