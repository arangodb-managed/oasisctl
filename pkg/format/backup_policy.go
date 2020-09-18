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
	backup "github.com/arangodb-managed/apis/backup/v1"
)

const (
	// Schedule Types
	hourly  = "Hourly"
	daily   = "Daily"
	monthly = "Monthly"
)

// BackupPolicy returns a single backup policy formatted for humans.
func BackupPolicy(x *backup.BackupPolicy, opts Options) string {
	data := backupPolicyToKeyValueList(x, opts)
	return formatObject(opts, data...)
}

func backupPolicyToKeyValueList(x *backup.BackupPolicy, opts Options) []kv {
	data := []kv{
		{"id", x.GetId()},
		{"deleted", formatBool(opts, x.GetIsDeleted())},
		{"deployment-id", x.GetDeploymentId()},
		{"description", x.GetDescription()},
		{"name", x.GetName()},
		{"upload", formatBool(opts, x.GetUpload())},
		{"url", x.GetUrl()},
		{"locked", formatBool(opts, x.GetLocked())},
		{"paused", formatBool(opts, x.GetIsPaused())},
		{"schedule-type", x.GetSchedule().GetScheduleType()},
		{"retention-period", formatDuration(opts, x.GetRetentionPeriod())},
		{"created-at", formatTime(opts, x.GetCreatedAt())},
		{"deleted-at", formatTime(opts, x.GetDeletedAt())},
	}

	if x.GetStatus() != nil {
		data = append(data, kv{"state", x.GetStatus().GetMessage()}, kv{"next-backup", formatTime(opts, x.GetStatus().GetNextBackup())})
	}

	switch x.GetSchedule().GetScheduleType() {
	case hourly:
		data = append(data, kv{
			"schedule-every-interval-hours",
			x.GetSchedule().GetHourlySchedule().GetScheduleEveryIntervalHours(),
		})
	case daily:
		dailySchedule := []kv{
			{"monday", formatBool(opts, x.GetSchedule().GetDailySchedule().GetMonday())},
			{"tuesday", formatBool(opts, x.GetSchedule().GetDailySchedule().GetTuesday())},
			{"wednesday", formatBool(opts, x.GetSchedule().GetDailySchedule().GetWednesday())},
			{"thursday", formatBool(opts, x.GetSchedule().GetDailySchedule().GetThursday())},
			{"friday", formatBool(opts, x.GetSchedule().GetDailySchedule().GetFriday())},
			{"saturday", formatBool(opts, x.GetSchedule().GetDailySchedule().GetSaturday())},
			{"sunday", formatBool(opts, x.GetSchedule().GetDailySchedule().GetSunday())},
			{"hour", x.GetSchedule().GetDailySchedule().GetScheduleAt().GetHours()},
			{"minutes", x.GetSchedule().GetDailySchedule().GetScheduleAt().GetMinutes()},
			{"timezone", x.GetSchedule().GetDailySchedule().GetScheduleAt().GetTimeZone()},
		}
		data = append(data, dailySchedule...)
	case monthly:
		monthlySchedule := []kv{
			{"day-of-month", x.GetSchedule().GetMonthlySchedule().GetDayOfMonth()},
			{"hour", x.GetSchedule().GetMonthlySchedule().GetScheduleAt().GetHours()},
			{"minutes", x.GetSchedule().GetMonthlySchedule().GetScheduleAt().GetMinutes()},
			{"timezone", x.GetSchedule().GetMonthlySchedule().GetScheduleAt().GetTimeZone()},
		}
		data = append(data, monthlySchedule...)
	}
	return data
}
