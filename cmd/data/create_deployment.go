//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package data

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	crypto "github.com/arangodb-managed/apis/crypto/v1"
	data "github.com/arangodb-managed/apis/data/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/format"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.CreateCmd,
		&cobra.Command{
			Use:   "deployment",
			Short: "Create a new deployment",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				name            string
				description     string
				organizationID  string
				projectID       string
				regionID        string
				cacertificateID string
				ipwhitelistID   string
				version         string
				serversPreset   string
				model           string
				nodeSizeId      string
				nodeCount       int32
				nodeDiskSize    int32
				// TODO add other fields
			}{}
			f.StringVar(&cargs.name, "name", "", "Name of the deployment")
			f.StringVar(&cargs.description, "description", "", "Description of the deployment")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization to create the deployment in")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project to create the deployment in")
			f.StringVarP(&cargs.regionID, "region-id", "r", cmd.DefaultRegion(), "Identifier of the region to create the deployment in")
			f.StringVarP(&cargs.cacertificateID, "cacertificate-id", "c", cmd.DefaultCACertificate(), "Identifier of the CA certificate to use for the deployment")
			f.StringVarP(&cargs.ipwhitelistID, "ipwhitelist-id", "i", cmd.DefaultIPWhitelist(), "Identifier of the IP whitelist to use for the deployment")
			f.StringVar(&cargs.version, "version", "", "Version of ArangoDB to use for the deployment")
			f.StringVar(&cargs.serversPreset, "servers-preset", "", "Servers preset to use for the deployment")
			f.StringVar(&cargs.model, "model", "oneshard", "Set model of the deployment")
			f.StringVar(&cargs.nodeSizeId, "node-size-id", "a4", "Set the node size to use for this deployment")
			f.Int32Var(&cargs.nodeCount, "node-count", 1, "Set the number of desired nodes")
			f.Int32Var(&cargs.nodeDiskSize, "node-disk-size", 20, "Set disk size for nodes (GB)")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				name, argsUsed := cmd.ReqOption("name", cargs.name, args, 0)
				regionID, _ := cmd.ReqOption("region-id", cargs.regionID, nil, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				cryptoc := crypto.NewCryptoServiceClient(conn)
				datac := data.NewDataServiceClient(conn)
				rmc := rm.NewResourceManagerServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Fetch project
				project := selection.MustSelectProject(ctx, log, cargs.projectID, cargs.organizationID, rmc)

				// Select cacertificate (to use in deployment)
				cacert := selection.MustSelectCACertificate(ctx, log, cargs.cacertificateID, project.GetId(), project.GetOrganizationId(), cryptoc, rmc)

				// Select servers from preset, if specified
				var servers *data.Deployment_ServersSpec
				if cargs.serversPreset != "" {
					servers = selection.MustSelectServersSpec(ctx, log, cargs.serversPreset, project.GetId(), regionID, datac)
				}

				// Create deployment
				result, err := datac.CreateDeployment(ctx, &data.Deployment{
					ProjectId:   project.GetId(),
					Name:        name,
					Description: cargs.description,
					RegionId:    regionID,
					Version:     cargs.version,
					Certificates: &data.Deployment_CertificateSpec{
						CaCertificateId: cacert.GetId(),
					},
					IpwhitelistId: cargs.ipwhitelistID,
					Servers:       servers,
					Model:         &data.Deployment_ModelSpec{
						Model:                cargs.model,
						NodeSizeId:           cargs.nodeSizeId,
						NodeCount:            cargs.nodeCount,
						NodeDiskSize:         cargs.nodeDiskSize,
					},
				})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to create deployment")
				}

				// Show result
				fmt.Println("Success!")
				fmt.Println(format.Deployment(result, nil, cmd.RootArgs.Format, false))
			}
		},
	)
}
