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

package selection

import (
	"context"

	"github.com/rs/zerolog"

	audit "github.com/arangodb-managed/apis/audit/v1"
	common "github.com/arangodb-managed/apis/common/v1"
)

// MustSelectAuditLogArchive fetches an auditlog archive given ID, name, or URL and fails if no such item is found.
func MustSelectAuditLogArchive(ctx context.Context, log zerolog.Logger, id string, auditc audit.AuditServiceClient) *audit.AuditLogArchive {
	result, err := SelectAuditLogArchive(ctx, log, id, auditc)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to select audit log archive")
	}
	return result
}

// SelectAuditLogArchive fetches an auditlog archive with given ID, name, or URL or returns an error if not found.
func SelectAuditLogArchive(ctx context.Context, log zerolog.Logger, id string, auditc audit.AuditServiceClient) (*audit.AuditLogArchive, error) {
	result, err := auditc.GetAuditLogArchive(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) {
			// Try to lookup auditlog by name or URL
			list, err := auditc.ListAuditLogArchives(ctx, &audit.ListAuditLogArchivesRequest{
				AuditlogId: id,
			})
			if err == nil {
				for _, x := range list.Items {
					if x.GetUrl() == id {
						return x, nil
					}
				}
			}
		}
		log.Debug().Err(err).Str("auditlog_archive_id", id).Msg("Failed to get audit log archive")
		return nil, err
	}
	return result, nil
}
