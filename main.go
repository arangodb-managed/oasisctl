//
// DISCLAIMER
//
// Copyright 2019 ArangoDB GmbH, Cologne, Germany
//
// Author Ewout Prangsma
//

package main

import (
	"log"

	_ "github.com/gogo/protobuf/types"

	"github.com/arangodb-managed/oasisctl/cmd"
	_ "github.com/arangodb-managed/oasisctl/cmd/crypto"
	_ "github.com/arangodb-managed/oasisctl/cmd/data"
	_ "github.com/arangodb-managed/oasisctl/cmd/iam"
	_ "github.com/arangodb-managed/oasisctl/cmd/platform"
	_ "github.com/arangodb-managed/oasisctl/cmd/resourcemanager"
	_ "github.com/arangodb-managed/oasisctl/cmd/security"
)

func init() {
	cmd.SetVersion(releaseVersion)
}

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalf("%v\n", err)
	}
}
