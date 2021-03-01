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
	// GetCmd is root for various `get ...` commands
	GetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get information",
		Run:   ShowUsage,
	}
	// GetMetricsCmd is root for various `get metrics ...` commands
	GetMetricsCmd = &cobra.Command{
		Use:   "metrics",
		Short: "Get metrics information",
		Run:   ShowUsage,
	}
)

func init() {
	RootCmd.AddCommand(GetCmd)
	GetCmd.AddCommand(GetMetricsCmd)
}
