//
// DISCLAIMER
//
// Copyright 2022 ArangoDB GmbH, Cologne, Germany
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

package format

import (
	"fmt"

	notebook "github.com/arangodb-managed/apis/notebook/v1"
)

// Notebook returns a single notebook formatted for humans.
func Notebook(x *notebook.Notebook, opts Options) string {
	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"deployment-id", x.GetDeploymentId()},
		kv{"name", x.GetName()},
		kv{"description", x.GetDescription()},
		kv{"url", x.GetUrl()},
		kv{"paused", formatBool(opts, x.GetIsPaused())},
		kv{"created-at", formatTime(opts, x.GetCreatedAt())},
		kv{"deleted-at", formatTime(opts, x.GetDeletedAt(), "-")},
		kv{"phase", x.GetStatus().GetPhase()},
		kv{"endpoint", x.GetStatus().GetEndpoint()},
	)
}

// NotebookList returns a list of notebooks formatted for humans.
func NotebookList(list []*notebook.Notebook, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			{"id", x.GetId()},
			{"name", x.GetName()},
			{"description", x.GetDescription()},
			{"url", x.GetUrl()},
			{"endpoint", x.GetStatus().GetEndpoint()},
			{"paused", formatBool(opts, x.GetIsPaused())},
			{"created-at", formatTime(opts, x.GetCreatedAt())},
		}
	}, false)
}

// NotebookModel returns a single NotebookModel formatted for humans.
func NotebookModel(x *notebook.NotebookModel, opts Options) string {
	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"name", x.GetName()},
		kv{"cpu", fmt.Sprintf("%0.2f", x.GetCpu())},
		kv{"memory", fmt.Sprintf("%dGiB", x.GetMemory())},
	)
}

// NotebookModelList returns a single NotebookModel formatted for humans.
func NotebookModelList(list []*notebook.NotebookModel, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			{"id", x.GetId()},
			{"name", x.GetName()},
			{"cpu", fmt.Sprintf("%.2f", x.GetCpu())},
			{"memory", fmt.Sprintf("%0.2dGiB", x.GetMemory())},
			{"minimum-disk-size", fmt.Sprintf("%0.2d GiB", x.GetMinDiskSize())},
			{"maximum-disk-size", fmt.Sprintf("%0.2d GiB", x.GetMaxDiskSize())},
		}
	}, false)
}
