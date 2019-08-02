//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package format

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gogo/protobuf/types"
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
func formatTime(opts Options, x *types.Timestamp, nilValue ...string) string {
	if x == nil {
		if len(nilValue) > 0 {
			return nilValue[0]
		}
		return ""
	}
	t, _ := types.TimestampFromProto(x)
	if opts.Format == formatJSON {
		return t.Format(time.RFC3339)
	}
	return humanize.Time(t)
}

// formatDuration returns a human readable version of the given duration.
func formatDuration(opts Options, x *types.Duration, nilValue ...string) string {
	if x == nil {
		if len(nilValue) > 0 {
			return nilValue[0]
		}
		return ""
	}
	d, _ := types.DurationFromProto(x)
	return d.String()
}
