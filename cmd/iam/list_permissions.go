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
