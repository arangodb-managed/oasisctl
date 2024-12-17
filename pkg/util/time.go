//
// DISCLAIMER
//
// Copyright 2020-2024 ArangoDB GmbH, Cologne, Germany
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

package util

import (
	"time"

	"github.com/araddon/dateparse"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ParseTimeFromNow parse a timestamp or duration before now.
func ParseTimeFromNow(value string) (time.Time, error) {
	if d, err := time.ParseDuration(value); err == nil {
		return time.Now().UTC().Add(-d), nil
	}
	ts, err := dateparse.ParseAny(value)
	if err != nil {
		return time.Time{}, err
	}
	return ts, nil
}

// ParseTime parses a given date string in RFC3339 to a proto timestamp.
// Usually used by from / to settings.
func ParseTime(date string) (*timestamppb.Timestamp, error) {
	d, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return nil, err
	}
	stamp := timestamppb.New(d)
	if err := stamp.CheckValid(); err != nil {
		return nil, err
	}
	return stamp, nil
}
