//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"

	"github.com/arangodb-managed/oasis/pkg/format"
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
	getCmd.AddCommand(getPolicyCmd)
	f := getPolicyCmd.Flags()
	f.StringVarP(&getPolicyArgs.url, "url", "u", defaultURL(), "URL of the resource to inspect the policy for")
}

func getPolicyCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	url, argsUsed := optOption("url", getPolicyArgs.url, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	ctx := contextWithToken()

	// Fetch policy
	item, err := iamc.GetPolicy(ctx, &common.URLOptions{Url: url})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to get policy")
	}

	// Show result
	fmt.Println(format.Policy(ctx, item, iamc, rootArgs.format))
}
