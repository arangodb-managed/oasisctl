//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package main

import (
	"log"

	_ "github.com/gogo/protobuf/types"

	"github.com/arangodb-managed/oasis/cmd"
	_ "github.com/arangodb-managed/oasis/cmd/crypto"
	_ "github.com/arangodb-managed/oasis/cmd/data"
	_ "github.com/arangodb-managed/oasis/cmd/iam"
	_ "github.com/arangodb-managed/oasis/cmd/platform"
	_ "github.com/arangodb-managed/oasis/cmd/resourcemanager"
	_ "github.com/arangodb-managed/oasis/cmd/security"
)

func init() {
	cmd.SetVersion(releaseVersion)
}

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalf("%v\n", err)
	}
}
