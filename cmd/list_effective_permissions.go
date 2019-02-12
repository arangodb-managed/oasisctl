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
	listEffectiveCmd.AddCommand(listEffectivePermissionsCmd)
	f := listEffectivePermissionsCmd.Flags()
	f.StringVarP(&listEffectivePermissionsArgs.url, "url", "u", defaultURL(), "URL of resource to get effective permissions for")
}

func listEffectivePermissionsCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	url, argsUsed := reqOption("url", listEffectivePermissionsArgs.url, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	ctx := contextWithToken()

	// Fetch permissions
	list, err := iamc.GetEffectivePermissions(ctx, &common.URLOptions{Url: url})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to list organizations")
	}

	// Show result
	fmt.Println(format.PermissionList(list.Items, rootArgs.format))
}
