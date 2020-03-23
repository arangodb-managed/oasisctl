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
	// getPolicyCmd fetches a project that the user has access to
	getPolicyCmd = &cobra.Command{
		Use:   "policy",
		Short: "Get a policy the authenticated user has access to",
		Run:   getPolicyCmdRun,
	}
	getPolicyArgs struct {
		url string
	}
)

func init() {
	cmd.GetCmd.AddCommand(getPolicyCmd)
	f := getPolicyCmd.Flags()
	f.StringVarP(&getPolicyArgs.url, "url", "u", cmd.DefaultURL(), "URL of the resource to inspect the policy for")
}

func getPolicyCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := getPolicyArgs
	url, argsUsed := cmd.OptOption("url", cargs.url, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch policy
	item, err := iamc.GetPolicy(ctx, &common.URLOptions{Url: url})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get policy")
	}

	// Show result
	fmt.Println(format.Policy(ctx, item, iamc, cmd.RootArgs.Format))
}
