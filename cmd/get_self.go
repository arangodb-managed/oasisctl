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
	// getSelfCmd fetches the user itself
	getSelfCmd = &cobra.Command{
		Use:   "self",
		Short: "Get information about the authenticated user",
		Run:   getSelfCmdRun,
	}
)

func init() {
	getCmd.AddCommand(getSelfCmd)
}

func getSelfCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	mustCheckNumberOfArgs(args, 0)

	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	ctx := contextWithToken()

	// Fetch user info
	user, err := iamc.GetThisUser(ctx, &common.Empty{})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to get user info")
	}

	// Show result
	fmt.Println(format.User(user, rootArgs.format))
}
