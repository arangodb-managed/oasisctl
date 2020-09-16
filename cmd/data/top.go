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
// Author Ewout Prangsma
//

package data

import (
	"fmt"
	"io"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gizak/termui/v3"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	common "github.com/arangodb-managed/apis/common/v1"
	data "github.com/arangodb-managed/apis/data/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.RootCmd,
		&cobra.Command{
			Use:   "top",
			Short: "Show the most important server metrics",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				deploymentID   string
				organizationID string
				projectID      string
			}{}
			f.StringVarP(&cargs.deploymentID, "deployment-id", "d", cmd.DefaultDeployment(), "Identifier of the deployment")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				deploymentID, argsUsed := cmd.OptOption("deployment-id", cargs.deploymentID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI(true)
				datac := data.NewDataServiceClient(conn)
				rmc := rm.NewResourceManagerServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Fetch deployment
				item := selection.MustSelectDeployment(ctx, log, deploymentID, cargs.projectID, cargs.organizationID, datac, rmc)

				// Prepare UI
				if err := ui.Init(); err != nil {
					log.Fatal().Err(err).Msg("failed to initialize UI")
				}
				defer ui.Close()
				w, h := termui.TerminalDimensions()
				// Create status
				sb := widgets.NewParagraph()
				sb.Text = ""
				sb.Border = false
				sb.SetRect(0, 0, w, 1)
				sb.TextStyle.Fg = ui.ColorMagenta
				// Create list
				list := widgets.NewList()
				list.TextStyle = ui.NewStyle(ui.ColorWhite)
				list.Border = false
				list.SetRect(0, 1, w, h-1)

				// Prepare updates
				updates := make(chan *data.Deployment)
				errors := make(chan error)

				// Start fetching updates
				go func() {
					defer func() {
						close(updates)
						close(errors)
					}()
					for {
						server, err := datac.GetDeploymentUpdates(ctx, &common.IDOptions{Id: item.GetId()})
						if err != nil {
							errors <- err
							return
						}
						for {
							depl, err := server.Recv()
							if err == io.EOF {
								// Connection closed normally, retry connection
								break
							} else if common.IsUnauthenticated(err) {
								errors <- err
								return
							} else if err != nil {
								errors <- err
								break
							} else {
								updates <- depl
							}
						}
					}
				}()

				// Render update & fetch events
				uiEvents := ui.PollEvents()
				lastUpdate := time.Now()
				depl := item
				var lastError error
				var lastErrorTS time.Time
				ok := true
				for {
					select {
					case depl, ok = <-updates:
						lastUpdate = time.Now()
						lastError = nil
					case lastError, ok = <-errors:
						lastErrorTS = time.Now()
					case e := <-uiEvents:
						switch e.ID {
						case "q", "<C-c>":
							return
						}
					case <-time.After(time.Second):
					}
					if !ok {
						return
					}
					if lastError != nil && time.Since(lastErrorTS) > time.Second*15 {
						lastError = nil
					}
					list.Rows = format.ServerStatusListAsRows(depl.GetStatus().GetServers(), cmd.RootArgs.Format)
					if lastError != nil {
						sb.TextStyle.Bg = ui.ColorRed
						sb.TextStyle.Fg = ui.ColorWhite
						sb.Text = fmt.Sprintf("Last update: %s, last error: %s", humanize.Time(lastUpdate), lastError)
					} else {
						sb.TextStyle.Bg = ui.ColorClear
						sb.TextStyle.Fg = ui.ColorWhite
						sb.Text = fmt.Sprintf("Last update: %s", humanize.Time(lastUpdate))
					}
					ui.Render(sb, list)
				}

			}
		},
	)
}
