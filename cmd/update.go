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

package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// UpdateCmd is root for various `update ...` commands
	UpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update resources",
		Run:   ShowUsage,
	}
	// UpdateMetricsCmd is root for various `update metrics ...` commands
	UpdateMetricsCmd = &cobra.Command{
		Use:   "metrics",
		Short: "Update metrics resources",
		Run:   ShowUsage,
	}
)

func init() {
	RootCmd.AddCommand(UpdateCmd)
	UpdateCmd.AddCommand(UpdateMetricsCmd)
}
