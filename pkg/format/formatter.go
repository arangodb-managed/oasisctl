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

package format

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/dustin/go-humanize"
	"github.com/ryanuber/columnize"
)

type kv struct {
	Key   string
	Value interface{}
}

var (
	singleConfig = &columnize.Config{
		Delim:  "|^|",
		Glue:   " ",
		Prefix: "",
		Empty:  "",
	}
	listConfig = &columnize.Config{
		Delim:  "|^|",
		Glue:   " | ",
		Prefix: "",
		Empty:  "",
	}
)

// formatObject returns a formatted representation of the given
// data which is a map from field-name to value.
func formatObject(opts Options, data ...kv) string {
	if opts.Format == formatJSON {
		m := make(map[string]interface{}, len(data))
		for _, kv := range data {
			m[kv.Key] = kv.Value
		}
		encoded, err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			panic(err)
		}
		return string(encoded)
	}

	// Table
	lines := make([]string, 0, len(data))
	for _, kv := range data {
		title := strings.Title(kv.Key)
		lines = append(lines, fmt.Sprintf("%s |^| %v", title, kv.Value))
	}
	return columnize.Format(lines, singleConfig)
}

// formatList returns a formatted representation of the given
// list.
func formatList(opts Options, list interface{}, getData func(int) []kv, noSort bool) string {
	listv := reflect.ValueOf(list)
	length := listv.Len()

	if opts.Format == formatJSON {
		l := make([]map[string]interface{}, length)
		for i := 0; i < length; i++ {
			data := getData(i)
			m := make(map[string]interface{})
			l[i] = m
			for _, kv := range data {
				m[kv.Key] = kv.Value
			}
		}
		encoded, err := json.MarshalIndent(l, "", "  ")
		if err != nil {
			panic(err)
		}
		return string(encoded)
	}

	// Table
	if length == 0 {
		return "None"
	}
	lines := make([]string, 0, length+2)
	for i := 0; i < length; i++ {
		data := getData(i)
		row := make([]string, 0, len(data))
		if i == 0 {
			for _, kv := range data {
				row = append(row, strings.Title(kv.Key))
			}
			lines = append(lines, strings.Join(row, "|^|"))
			row = row[:0]
		}
		for _, kv := range data {
			row = append(row, fmt.Sprintf("%v", kv.Value))
		}
		lines = append(lines, strings.Join(row, "|^|"))
	}
	if !noSort {
		sort.Strings(lines[1:])
	}
	return columnize.Format(lines, listConfig)
}

// formatTime returns a human readable version of the given timestamp.
func formatTime(opts Options, x *timestamppb.Timestamp, nilValue ...string) string {
	if x == nil {
		if len(nilValue) > 0 {
			return nilValue[0]
		}
		return ""
	}
	t := x.AsTime()
	if opts.Format == formatJSON {
		return t.Format(time.RFC3339)
	}
	return humanize.Time(t)
}

// formatDuration returns a human readable version of the given duration.
func formatDuration(opts Options, x *durationpb.Duration, nilValue ...string) string {
	if x == nil {
		if len(nilValue) > 0 {
			return nilValue[0]
		}
		return ""
	}
	d := x.AsDuration()
	return d.String()
}

// formatOptionalString returns the given string or "-" when input is empty.
func formatOptionalString(value string) string {
	if value == "" {
		return "-"
	}
	return value
}
