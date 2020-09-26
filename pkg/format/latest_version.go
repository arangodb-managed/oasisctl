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

package format

import (
	tools "github.com/arangodb-managed/apis/tools/v1"
)

// ToolsVersion returns a single tools version formatted for humans.
func ToolsVersion(x *tools.ToolsVersion, opts Options) string {
	return formatObject(opts,
		kv{"url", x.GetDownloadUrl()},
		kv{"latest-version", x.GetLatestVersion()},
	)
}
