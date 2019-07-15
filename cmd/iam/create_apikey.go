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
	flag "github.com/spf13/pflag"

	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/format"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.CreateCmd,
		&cobra.Command{
			Use:   "apikey",
			Short: "Create a new API key",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				readonly       bool
				organizationID string
			}{}
			f.BoolVar(&cargs.readonly, "readonly", false, "If set, the newly created API key will grant readonly access only")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", "", "If set, the newly created API key will grant access to this organization only")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				cmd.MustCheckNumberOfArgs(args, 0)

				// Connect
				conn := cmd.MustDialAPI()
				iamc := iam.NewIAMServiceClient(conn)
				rmc := rm.NewResourceManagerServiceClient(conn)
				ctx := cmd.ContextWithToken()

				var orgID string
				// Fetch organization
				if cargs.organizationID != "" {
					org := selection.MustSelectOrganization(ctx, log, cargs.organizationID, rmc)
					orgID = org.GetId()
				}

				// Create API key
				result, err := iamc.CreateAPIKey(ctx, &iam.CreateAPIKeyRequest{
					Readonly:       cargs.readonly,
					OrganizationId: orgID,
				})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to create API key")
				}

				// Show result
				fmt.Println("Success!")
				fmt.Println(format.APIKeySecret(result, cmd.RootArgs.Format))
			}
		},
	)
}
