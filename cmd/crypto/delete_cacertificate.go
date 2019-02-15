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

	common "github.com/arangodb-managed/apis/common/v1"
	crypto "github.com/arangodb-managed/apis/crypto/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

var (
	// deleteCACertificateCmd deletes a role that the user has access to
	deleteCACertificateCmd = &cobra.Command{
		Use:   "cacertificate",
		Short: "Delete a CA certificate the authenticated user has access to",
		Run:   deleteCACertificateCmdRun,
	}
	deleteCACertificateArgs struct {
		organizationID string
		projectID      string
		cacertID       string
	}
)

func init() {
	cmd.DeleteCmd.AddCommand(deleteCACertificateCmd)
	f := deleteCACertificateCmd.Flags()
	f.StringVarP(&deleteCACertificateArgs.cacertID, "cacertificate-id", "c", cmd.DefaultCACertificate(), "Identifier of the CA certificate")
	f.StringVarP(&deleteCACertificateArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
	f.StringVarP(&deleteCACertificateArgs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")
}

func deleteCACertificateCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	cacertID, argsUsed := cmd.OptOption("cacertificate-id", deleteCACertificateArgs.cacertID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	cryptoc := crypto.NewCryptoServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch CA certificate
	item := selection.MustSelectCACertificate(ctx, cmd.CLILog, cacertID, deleteCACertificateArgs.projectID, deleteCACertificateArgs.organizationID, cryptoc, rmc)

	// Delete CA certificate
	if _, err := cryptoc.DeleteCACertificate(ctx, &common.IDOptions{Id: item.GetId()}); err != nil {
		cmd.CLILog.Fatal().Err(err).Msg("Failed to delete CA certificate")
	}

	// Show result
	fmt.Println("Deleted CA certificate!")
}
