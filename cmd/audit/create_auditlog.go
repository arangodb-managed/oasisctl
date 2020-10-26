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

package audit

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	audit "github.com/arangodb-managed/apis/audit/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.CreateCmd,
		&cobra.Command{
			Use:   "auditlog",
			Short: "Create an auditlog",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				name                 string
				description          string
				organizationID       string
				isDefault            bool
				destinationType      string
				url                  string
				trustedServerCAPem   string
				clientCertificatePem string
				clientKeyPem         string
				headers              []string
				excludedTopics       []string
			}{}
			f.StringVar(&cargs.name, "name", "", "Name of the audit log.")
			f.StringVar(&cargs.description, "description", "", "Description of the audit log.")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
			f.BoolVar(&cargs.isDefault, "default", false, "If set, this AuditLog is the default for the organization.")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				orgID, argsUsed := cmd.OptOption("organization-id", cargs.organizationID, args, 0)
				name, argsUsed := cmd.OptOption("name", cargs.name, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				auditc := audit.NewAuditServiceClient(conn)
				rmc := rm.NewResourceManagerServiceClient(conn)
				ctx := cmd.ContextWithToken()
				org := selection.MustSelectOrganization(ctx, log, orgID, rmc)

				// Construct the default destination
				destination := &audit.AuditLog_Destination{
					Type: audit.DestinationCloud,
				}

				// Construct request
				req := &audit.AuditLog{
					Name:           name,
					Description:    cargs.description,
					OrganizationId: org.GetId(),
					IsDefault:      cargs.isDefault,
					Destinations:   []*audit.AuditLog_Destination{destination},
				}

				// Make the call
				result, err := auditc.CreateAuditLog(ctx, req)
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to create audit log.")
				}

				// Show result
				format.DisplaySuccess(cmd.RootArgs.Format)
				fmt.Println(format.AuditLog(result, cmd.RootArgs.Format))
			}
		},
	)
}
