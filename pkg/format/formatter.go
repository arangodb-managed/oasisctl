//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package format

import (
	"fmt"
	"reflect"
	"strings"

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
	lines := make([]string, 0, len(data))
	for _, kv := range data {
		title := strings.Title(kv.Key)
		lines = append(lines, fmt.Sprintf("%s |^| %v", title, kv.Value))
	}
	return columnize.Format(lines, singleConfig)
}

// formatList returns a formatted representation of the given
// list.
func formatList(opts Options, list interface{}, getData func(int) []kv) string {
	listv := reflect.ValueOf(list)
	length := listv.Len()
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
	return columnize.Format(lines, listConfig)
}

// formatTime returns a human readable version of the given timestamp.
func formatTime(x *types.Timestamp, nilValue ...string) string {
	if x == nil {
		if len(nilValue) > 0 {
			return nilValue[0]
		}
		return ""
	}
	t, _ := types.TimestampFromProto(x)
	return humanize.Time(t)
}
