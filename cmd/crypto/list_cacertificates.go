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
	"github.com/arangodb-managed/oasis/pkg/format"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

var (
	// listCACertificatesCmd fetches CA certificates of the given project
	listCACertificatesCmd = &cobra.Command{
		Use:   "cacertificates",
		Short: "List all CA certificates of the given project",
		Run:   listCACertificatesCmdRun,
	}
	listCACertificatesArgs struct {
		organizationID string
		projectID      string
	}
)

func init() {
	cmd.ListCmd.AddCommand(listCACertificatesCmd)
	f := listCACertificatesCmd.Flags()
	f.StringVarP(&listCACertificatesArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
	f.StringVarP(&listCACertificatesArgs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")
}

func listCACertificatesCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := listCACertificatesArgs
	projectID, argsUsed := cmd.OptOption("project-id", cargs.projectID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	cryptoc := crypto.NewCryptoServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch project
	project := selection.MustSelectProject(ctx, log, projectID, cargs.organizationID, rmc)

	// Fetch CA certificates in project
	list, err := cryptoc.ListCACertificates(ctx, &common.ListOptions{ContextId: project.GetId()})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list CA certificates")
	}

	// Show result
	fmt.Println(format.CACertificateList(list.Items, cmd.RootArgs.Format))
}
