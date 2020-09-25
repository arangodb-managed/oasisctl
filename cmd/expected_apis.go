//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Author Gergely Brautigam
//

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	backup "github.com/arangodb-managed/apis/backup/v1"
	crypto "github.com/arangodb-managed/apis/crypto/v1"
	data "github.com/arangodb-managed/apis/data/v1"
	example "github.com/arangodb-managed/apis/example/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"
	mon "github.com/arangodb-managed/apis/monitoring/v1"
	platform "github.com/arangodb-managed/apis/platform/v1"
	replication "github.com/arangodb-managed/apis/replication/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
	security "github.com/arangodb-managed/apis/security/v1"
)

const (
	jsonFilename = "apis.json"
)

func init() {
	InitCommand(
		RootCmd,
		&cobra.Command{
			Use:    "expected-apis",
			Short:  "Generate an apis.json file.",
			Long:   "Generates a file which contains all the versions needed by this tool.",
			Hidden: true,
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			c.Run = func(c *cobra.Command, args []string) {
				convertVersionToString := func(major int, minor int, patch int) string {
					return fmt.Sprintf("%d.%d.%d", major, minor, patch)
				}
				versions := make(map[string]string)
				versions[iam.APIID] = convertVersionToString(iam.APIMajorVersion, iam.APIMinorVersion, iam.APIPatchVersion)
				versions[backup.APIID] = convertVersionToString(backup.APIMajorVersion, backup.APIMinorVersion, backup.APIPatchVersion)
				versions[crypto.APIID] = convertVersionToString(crypto.APIMajorVersion, crypto.APIMinorVersion, crypto.APIPatchVersion)
				versions[data.APIID] = convertVersionToString(data.APIMajorVersion, data.APIMinorVersion, data.APIPatchVersion)
				versions[example.APIID] = convertVersionToString(example.APIMajorVersion, example.APIMinorVersion, example.APIPatchVersion)
				versions[mon.APIID] = convertVersionToString(mon.APIMajorVersion, mon.APIMinorVersion, mon.APIPatchVersion)
				versions[platform.APIID] = convertVersionToString(platform.APIMajorVersion, platform.APIMinorVersion, platform.APIPatchVersion)
				versions[replication.APIID] = convertVersionToString(replication.APIMajorVersion, replication.APIMinorVersion, replication.APIPatchVersion)
				versions[rm.APIID] = convertVersionToString(rm.APIMajorVersion, rm.APIMinorVersion, rm.APIPatchVersion)
				versions[security.APIID] = convertVersionToString(security.APIMajorVersion, security.APIMinorVersion, security.APIPatchVersion)
				content, err := json.Marshal(versions)
				if err != nil {
					CLILog.Fatal().Err(err).Msg("Failed to marshal map to json.")
				}
				if err := ioutil.WriteFile(jsonFilename, content, 0755); err != nil {
					CLILog.Fatal().Err(err).Msg("Failed to write out file.")
				}
			}
		},
	)
}
