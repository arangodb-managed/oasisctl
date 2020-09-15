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

package format

import (
	crypto "github.com/arangodb-managed/apis/crypto/v1"
)

// CACertificate returns a single ca certificate formatted for humans.
func CACertificate(x *crypto.CACertificate, opts Options) string {
	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"name", x.GetName()},
		kv{"description", x.GetDescription()},
		kv{"lifetime", formatDuration(opts, x.GetLifetime())},
		kv{"url", x.GetUrl()},
		kv{"use-well-known-certificate", formatBool(opts, x.GetUseWellKnownCertificate())},
		kv{"locked", formatBool(opts, x.GetLocked())},
		kv{"created-at", formatTime(opts, x.GetCreatedAt())},
		kv{"deleted-at", formatTime(opts, x.GetDeletedAt(), "-")},
	)
}

// CACertificateList returns a list of ca certificates formatted for humans.
func CACertificateList(list []*crypto.CACertificate, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			{"id", x.GetId()},
			{"name", x.GetName()},
			{"description", x.GetDescription()},
			{"lifetime", formatDuration(opts, x.GetLifetime())},
			{"url", x.GetUrl()},
			{"use-well-known-certificate", formatBool(opts, x.GetUseWellKnownCertificate())},
			{"locked", formatBool(opts, x.GetLocked())},
			{"created-at", formatTime(opts, x.GetCreatedAt())},
		}
	}, false)
}
