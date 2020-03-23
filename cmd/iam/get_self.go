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

package iam

import (
	"fmt"

	"github.com/spf13/cobra"

	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
)

var (
	// getSelfCmd fetches the user itself
	getSelfCmd = &cobra.Command{
		Use:   "self",
		Short: "Get information about the authenticated user",
		Run:   getSelfCmdRun,
	}
)

func init() {
	cmd.GetCmd.AddCommand(getSelfCmd)
}

func getSelfCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cmd.MustCheckNumberOfArgs(args, 0)

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch user info
	user, err := iamc.GetThisUser(ctx, &common.Empty{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get user info")
	}

	// Show result
	fmt.Println(format.User(user, cmd.RootArgs.Format))
}
