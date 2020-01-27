//
// DISCLAIMER
//
// Copyright 2019 ArangoDB GmbH, Cologne, Germany
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
