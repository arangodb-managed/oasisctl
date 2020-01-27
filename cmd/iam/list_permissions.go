//
// DISCLAIMER
//
// Copyright 2019 ArangoDB GmbH, Cologne, Germany
//
// Author Robert Stam
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
	// listPermissionsCmd fetches the known permissions.
	listPermissionsCmd = &cobra.Command{
		Use:   "permissions",
		Short: "List the known permissions",
		Run:   listPermissionsCmdRun,
	}
	listPermissionsArgs struct {
	}
)

func init() {
	cmd.ListCmd.AddCommand(listPermissionsCmd)
}

func listPermissionsCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch permissions
	list, err := iamc.ListPermissions(ctx, &common.Empty{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list permissions")
	}

	// Show result
	fmt.Println(format.PermissionList(list.Items, cmd.RootArgs.Format))
}
