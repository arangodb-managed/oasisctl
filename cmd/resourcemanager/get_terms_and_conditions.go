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
// Author Gergely Brautigam
//

package rm

import (
	"fmt"

	"github.com/spf13/cobra"

	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

var (
	// getTermsAndConditionsCmd fetches a tandc that the user has access to
	getTermsAndConditionsCmd = &cobra.Command{
		Use:   "tandc",
		Short: "Get current terms and conditions or get one by ID",
		Run:   getTermsAndConditionsCmdRun,
	}
	getTermsAndConditionsArgs struct {
		organizationID string
		tandcID        string
	}
)

func init() {
	cmd.GetCmd.AddCommand(getTermsAndConditionsCmd)
	f := getTermsAndConditionsCmd.Flags()
	f.StringVarP(&getTermsAndConditionsArgs.tandcID, "terms-and-conditions-id", "t", cmd.DefaultTermsAndConditions(), "Identifier of the terms and conditions to accept.")
	f.StringVarP(&getTermsAndConditionsArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
}

func getTermsAndConditionsCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := getTermsAndConditionsArgs
	tandcID, argsUsed := cmd.OptOption("terms-and-conditions-id", cargs.tandcID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch tandc
	item := selection.MustSelectTermsAndConditions(ctx, log, tandcID, cargs.organizationID, rmc)

	// Show result
	fmt.Println(format.TermsAndConditions(item, cmd.RootArgs.Format))
}
