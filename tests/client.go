package tests

import (
	"context"
	"os"

	"github.com/rs/zerolog"

	crypto "github.com/arangodb-managed/apis/crypto/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

// GetCryptoClientAndProject creates a crypto client and a project for the tests to work with.
func GetCryptoClientAndProject(ctx context.Context) (crypto.CryptoServiceClient, *rm.Project) {
	log := zerolog.New(zerolog.ConsoleWriter{
		Out:     os.Stderr,
		NoColor: true,
	}).With().Timestamp().Logger()
	org := cmd.DefaultOrganization()
	proj := cmd.DefaultProject()

	conn := cmd.MustDialAPI()
	cryptoc := crypto.NewCryptoServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	project := selection.MustSelectProject(ctx, log, proj, org, rmc)
	return cryptoc, project
}
