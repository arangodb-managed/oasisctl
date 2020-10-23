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
	"strings"

	"github.com/arangodb-managed/oasisctl/pkg/selection"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	audit "github.com/arangodb-managed/apis/audit/v1"
	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
)

const (
	cloud    = "cloud"
	httpPost = "https-post"
)

func init() {
	cmd.InitCommand(
		cmd.AddAuditLogCmd,
		&cobra.Command{
			Use:   "destination",
			Short: "Add a destination to an auditlog.",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				auditLogID           string
				destinationType      string
				url                  string
				trustedServerCAPem   string
				clientCertificatePem string
				clientKeyPem         string
				headers              []string
				excludedTopics       []string
				organizationID       string
			}{}
			f.StringVarP(&cargs.auditLogID, "auditlog-id", "i", "", "Identifier of the auditlog")
			f.StringVar(&cargs.destinationType, "destination-type", "", `Type of destination. Possible values are: "cloud", "https-post"`)
			f.StringVar(&cargs.url, "destination-https-url", "", "URL of the server to POST to. Scheme must be HTTPS.")
			f.StringVar(&cargs.trustedServerCAPem, "destination-https-trusted-server-ca-pem", "", "PEM encoded public key of the CA used to sign the server TLS certificate. If empty, a well known CA is expected.")
			f.StringVar(&cargs.clientCertificatePem, "destination-https-client-certificate-pem", "", "URL of the server to POST to. Scheme must be HTTPS.")
			f.StringVar(&cargs.clientKeyPem, "destination-https-client-key-pem", "", "URL of the server to POST to. Scheme must be HTTPS.")
			f.StringSliceVar(&cargs.headers, "destination-https-headers", nil, "A key=value formatted list of headers for the request. Repeating headers are allowed.")
			f.StringSliceVar(&cargs.excludedTopics, "destination-https-excluded-topics", nil, "Do not send audit events with these topics to this destination.")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				id, argsUsed := cmd.ReqOption("auditlog-id", cargs.auditLogID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				auditc := audit.NewAuditServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Construct destination
				destination := &audit.AuditLog_Destination{}
				switch cargs.destinationType {
				case cloud:
					destination.Type = cloud
				case httpPost:
					// Construct header values
					headers := make([]*audit.AuditLog_Header, 0)
					for _, h := range cargs.headers {
						split := strings.Split(h, "=")
						if len(split) < 2 {
							log.Fatal().Strs("headers", cargs.headers).Msg("Headers must be of format key=value")
						}
						if split[0] == "" {
							log.Fatal().Strs("headers", cargs.headers).Msg("Header name must not be empty.")
						}
						headers = append(headers, &audit.AuditLog_Header{
							Key:   split[0],
							Value: split[1],
						})
					}

					destination.Type = httpPost
					destination.HttpPost = &audit.AuditLog_HttpsPostSettings{
						Url:                  cargs.url,
						TrustedServerCaPem:   cargs.trustedServerCAPem,
						ClientCertificatePem: cargs.clientCertificatePem,
						ClientKeyPem:         cargs.clientKeyPem,
						Headers:              headers,
						ExcludedTopics:       cargs.excludedTopics,
					}
				default:
					log.Fatal().Str("type", cargs.destinationType).Msg(`Invalid destination type. Can be one of "cloud" or "https-post"`)
				}

				// Get the auditlog
				auditLog := selection.MustSelectAuditLog(ctx, log, id, cargs.organizationID, auditc)

				// Add the destination
				auditLog.Destinations = append(auditLog.GetDestinations(), destination)

				// Update auditlog with the new destinations.
				result, err := auditc.UpdateAuditLog(ctx, auditLog)
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to add new destination.")
				}

				// Show result
				format.DisplaySuccess(cmd.RootArgs.Format)
				fmt.Println(format.AuditLog(result, cmd.RootArgs.Format))
			}
		},
	)
}
