//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
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
	// listEffectivePermissionsCmd fetches the effective permissions of the user for a given URL.
	listEffectivePermissionsCmd = &cobra.Command{
		Use:   "permissions",
		Short: "List the effective permissions, the authenticated user has for a given URL",
		Run:   listEffectivePermissionsCmdRun,
	}
	listEffectivePermissionsArgs struct {
		url string
	}
)

func init() {
	cmd.ListEffectiveCmd.AddCommand(listEffectivePermissionsCmd)
	f := listEffectivePermissionsCmd.Flags()
	f.StringVarP(&listEffectivePermissionsArgs.url, "url", "u", cmd.DefaultURL(), "URL of resource to get effective permissions for")
}

func listEffectivePermissionsCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := listEffectivePermissionsArgs
	url, argsUsed := cmd.ReqOption("url", cargs.url, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch permissions
	list, err := iamc.GetEffectivePermissions(ctx, &common.URLOptions{Url: url})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list effective permissions")
	}

	// Show result
	fmt.Println(format.PermissionList(list.Items, cmd.RootArgs.Format))
}
