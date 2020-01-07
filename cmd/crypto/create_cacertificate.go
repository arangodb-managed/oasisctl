//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package crypto

import (
	"fmt"
	"time"

	types "github.com/gogo/protobuf/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	crypto "github.com/arangodb-managed/apis/crypto/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.CreateCmd,
		&cobra.Command{
			Use:   "cacertificate",
			Short: "Create a new CA certificate",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				name           string
				description    string
				organizationID string
				projectID      string
				lifetime       time.Duration
			}{}
			f.StringVar(&cargs.name, "name", "", "Name of the CA certificate")
			f.StringVar(&cargs.description, "description", "", "Description of the CA certificate")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization to create the CA certificate in")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project to create the CA certificate in")
			f.DurationVar(&cargs.lifetime, "lifetime", 0, "Lifetime of the CA certificate.")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				name, argsUsed := cmd.ReqOption("name", cargs.name, args, 0)
				description := cargs.description
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				cryptoc := crypto.NewCryptoServiceClient(conn)
				rmc := rm.NewResourceManagerServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Fetch project
				project := selection.MustSelectProject(ctx, log, cargs.projectID, cargs.organizationID, rmc)

				// Create ca certificate
				var lifetime *types.Duration
				if cargs.lifetime > 0 {
					lifetime = types.DurationProto(cargs.lifetime)
				}
				result, err := cryptoc.CreateCACertificate(ctx, &crypto.CACertificate{
					ProjectId:   project.GetId(),
					Name:        name,
					Description: description,
					Lifetime:    lifetime,
				})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to create CA certificate")
				}

				// Show result
				fmt.Println("Success!")
				fmt.Println(format.CACertificate(result, cmd.RootArgs.Format))
			}
		},
	)
}
