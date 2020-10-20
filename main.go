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
// Author Ewout Prangsma
//

package main

import (
	"log"
	"os"

	_ "github.com/gogo/protobuf/types"

	"github.com/arangodb-managed/oasisctl/cmd"
	_ "github.com/arangodb-managed/oasisctl/cmd/audit"
	_ "github.com/arangodb-managed/oasisctl/cmd/crypto"
	_ "github.com/arangodb-managed/oasisctl/cmd/data"
	_ "github.com/arangodb-managed/oasisctl/cmd/example"
	_ "github.com/arangodb-managed/oasisctl/cmd/iam"
	_ "github.com/arangodb-managed/oasisctl/cmd/importdata"
	_ "github.com/arangodb-managed/oasisctl/cmd/platform"
	_ "github.com/arangodb-managed/oasisctl/cmd/resourcemanager"
	_ "github.com/arangodb-managed/oasisctl/cmd/security"
)

func init() {
	cmd.SetVersion(releaseVersion)
}

func main() {
	var pluginHandler cmd.PluginHandler
	switch cmd.DefaultPluginHandler() {
	case "", "default":
		pluginHandler = cmd.NewDefaultPluginHandler("oasisctl")
	case "none":
		pluginHandler = nil
	default:
		log.Fatalf("Unknown plugin handler '%s'\n", cmd.DefaultPluginHandler())
	}
	cmd.ExecuteCommandOrPlugin(cmd.RootCmd, pluginHandler, os.Args)
}
