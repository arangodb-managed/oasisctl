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
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalf("%v\n", err)
	}
}
