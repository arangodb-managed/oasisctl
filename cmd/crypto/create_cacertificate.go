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

	crypto "github.com/arangodb-managed/apis/crypto/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/format"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

var (
	// createCACertificateCmd creates a new CA certificate
	createCACertificateCmd = &cobra.Command{
		Use:   "cacertificate",
		Short: "Create a new CA certificate",
		Run:   createCACertificateCmdRun,
	}
	createCACertificateArgs struct {
		name           string
		description    string
		organizationID string
		projectID      string
		lifetime       time.Duration
	}
)

func init() {
	cmd.CreateCmd.AddCommand(createCACertificateCmd)

	f := createCACertificateCmd.Flags()
	f.StringVarP(&createCACertificateArgs.name, "name", "n", "", "Name of the CA certificate")
	f.StringVarP(&createCACertificateArgs.description, "description", "d", "", "Description of the CA certificate")
	f.StringVarP(&createCACertificateArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization to create the CA certificate in")
	f.StringVarP(&createCACertificateArgs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project to create the CA certificate in")
	f.DurationVar(&createCACertificateArgs.lifetime, "lifetime", 0, "Lifetime of the CA certificate.")
}

func createCACertificateCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := createCACertificateArgs
	name, argsUsed := cmd.ReqOption("name", cargs.name, args, 0)
	description := createCACertificateArgs.description
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
