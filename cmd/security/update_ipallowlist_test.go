//
// DISCLAIMER
//
// Copyright 2020-2023 ArangoDB GmbH, Cologne, Germany
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
//

package security

import (
	"testing"

	security "github.com/arangodb-managed/apis/security/v1"
)

func TestUpdateCidrRanges(t *testing.T) {
	testCases := []struct {
		name               string
		existingCidrRanges []string
		cidrRanges         map[string]struct{}
		addCidrRanges      []string
		removeCidrRanges   []string
		expectedCidrRanges []string
	}{
		{
			name:               "Add CIDR range and ensure order remains the same",
			existingCidrRanges: []string{"203.0.113.0/24", "10.0.0.0/8"},
			cidrRanges: map[string]struct{}{
				"192.168.0.0/24": {},
				"10.0.0.0/8":     {},
			},
			addCidrRanges:      []string{"172.16.0.0/16", "88.11.0.0/16"},
			removeCidrRanges:   []string{},
			expectedCidrRanges: []string{"203.0.113.0/24", "10.0.0.0/8", "172.16.0.0/16", "88.11.0.0/16"},
		},
		{
			name:               "Remove CIDR range",
			existingCidrRanges: []string{"88.11.0.0/16", "10.0.0.0/8"},
			cidrRanges: map[string]struct{}{
				"88.11.0.0/16": {},
			},
			addCidrRanges:      []string{},
			removeCidrRanges:   []string{"10.0.0.0/8"},
			expectedCidrRanges: []string{"88.11.0.0/16"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			item := &security.IPAllowlist{
				CidrRanges: tc.existingCidrRanges,
			}

			UpdateCidrRanges(tc.existingCidrRanges, tc.cidrRanges, &UpdateCidrRangeOptions{
				AddedCidrRanges:   tc.addCidrRanges,
				RemovedCidrRanges: tc.removeCidrRanges,
			}, item)

			if !stringSlicesEqual(item.CidrRanges, tc.expectedCidrRanges) {
				t.Errorf("Unexpected CIDR ranges. Got %v, expected %v", item.CidrRanges, tc.expectedCidrRanges)
			}
		})
	}
}

// Helper function to compare string slices
func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
