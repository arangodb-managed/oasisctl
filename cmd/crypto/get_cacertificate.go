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

	"github.com/spf13/cobra"

	crypto "github.com/arangodb-managed/apis/crypto/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/format"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

var (
	// getCACertificateCmd fetches a CA certificate that the user has access to
	getCACertificateCmd = &cobra.Command{
		Use:   "cacertificate",
		Short: "Get a CA certificate the authenticated user has access to",
		Run:   getCACertificateCmdRun,
	}
	getCACertificateArgs struct {
		cacertID       string
		organizationID string
		projectID      string
	}
)

func init() {
	cmd.GetCmd.AddCommand(getCACertificateCmd)
	f := getCACertificateCmd.Flags()
	f.StringVarP(&getCACertificateArgs.cacertID, "cacertificate-id", "c", cmd.DefaultCACertificate(), "Identifier of the CA certificate")
	f.StringVarP(&getCACertificateArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
	f.StringVarP(&getCACertificateArgs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")
}

func getCACertificateCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	cacertID, argsUsed := cmd.OptOption("cacertificate-id", getCACertificateArgs.cacertID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	cryptoc := crypto.NewCryptoServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch CA certificate
	item := selection.MustSelectCACertificate(ctx, cmd.CLILog, cacertID, getCACertificateArgs.projectID, getCACertificateArgs.organizationID, cryptoc, rmc)

	// Show result
	fmt.Println(format.CACertificate(item, cmd.RootArgs.Format))
}
