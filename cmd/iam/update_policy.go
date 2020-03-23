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
	"github.com/spf13/cobra"

	"github.com/arangodb-managed/oasisctl/cmd"
)

var (
	// updatePolicyCmd is root for various `update policy ...` commands
	updatePolicyCmd = &cobra.Command{
		Use:   "policy",
		Short: "Update a policy",
		Run:   cmd.ShowUsage,
	}
	// updatePolicyAddCmd is root for various `update policy add ...` commands
	updatePolicyAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add to a policy",
		Run:   cmd.ShowUsage,
	}
	// updatePolicyDeleteCmd is root for various `update policy delete ...` commands
	updatePolicyDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete from a policy",
		Run:   cmd.ShowUsage,
	}
)

func init() {
	cmd.UpdateCmd.AddCommand(updatePolicyCmd)
	updatePolicyCmd.AddCommand(updatePolicyAddCmd)
	updatePolicyCmd.AddCommand(updatePolicyDeleteCmd)
}
