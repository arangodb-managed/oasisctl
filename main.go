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
	_ "github.com/arangodb-managed/oasis/cmd/iam"
	_ "github.com/arangodb-managed/oasis/cmd/platform"
	_ "github.com/arangodb-managed/oasis/cmd/resourcemanager"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalf("%v\n", err)
	}
}
