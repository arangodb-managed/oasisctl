//
// DISCLAIMER
//
// Copyright 2020-2022 ArangoDB GmbH, Cologne, Germany
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

package data

import (
	"fmt"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	backup "github.com/arangodb-managed/apis/backup/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

func init() {
	cmd.InitCommand(
		updateBackupCmd,
		&cobra.Command{
			Use:   "policy",
			Short: "Update a backup policy",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				id                string
				name              string
				description       string
				emailNotification string
				scheduleType      string
				paused            bool
				hourlySchedule    struct {
					scheduleEveryIntervalHours int32
					minutesOffset              int32
				}
				dailySchedule struct {
					monday    bool
					tuesday   bool
					wednesday bool
					thursday  bool
					friday    bool
					saturday  bool
					sunday    bool
				}
				monthlySchedule struct {
					dayOfMonth int32
				}
				timeofday struct {
					hours    int32
					minutes  int32
					timezone string
				}
				retentionPeriod     int
				upload              bool
				locked              bool
				additionalRegionIDs []string
			}{}
			f.StringVarP(&cargs.id, "backup-policy-id", "d", "", "Identifier of the backup policy")
			f.StringVar(&cargs.name, "name", "", "Name of the deployment")
			f.StringVar(&cargs.description, "description", "", "Description of the backup")
			f.StringVar(&cargs.emailNotification, "email-notification", "", "Email notification setting (Never|FailureOnly|Always)")
			f.StringVar(&cargs.scheduleType, "schedule-type", "", "Schedule of the policy (Hourly|Daily|Monthly)")
			f.BoolVar(&cargs.upload, "upload", false, "The backup should be uploaded. Set to false explicitly to clear the flag.")
			f.BoolVar(&cargs.paused, "paused", false, "The policy is paused. Set to false explicitly to clear the flag.")
			f.IntVar(&cargs.retentionPeriod, "retention-period", 0, "Backups created by this policy will be automatically deleted after the specified retention period. A value of 0 means that backup will never be deleted.")
			f.Int32Var(&cargs.hourlySchedule.scheduleEveryIntervalHours, "every-interval-hours", 0, "Schedule should run with an interval of the specified hours (1-23)")
			f.Int32Var(&cargs.hourlySchedule.minutesOffset, "minutes-offset", 0, "Schedule should run with specific minutes offset (0-59)")
			f.BoolVar(&cargs.dailySchedule.monday, "monday", false, "If set, a backup will be created on Mondays. Set to false explicitly to clear the flag.")
			f.BoolVar(&cargs.dailySchedule.tuesday, "tuesday", false, "If set, a backup will be created on Tuesdays. Set to false explicitly to clear the flag.")
			f.BoolVar(&cargs.dailySchedule.wednesday, "wednesday", false, "If set, a backup will be created on Wednesdays. Set to false explicitly to clear the flag.")
			f.BoolVar(&cargs.dailySchedule.thursday, "thursday", false, "If set, a backup will be created on Thursdays. Set to false explicitly to clear the flag.")
			f.BoolVar(&cargs.dailySchedule.friday, "friday", false, "If set, a backup will be created on Fridays. Set to false explicitly to clear the flag.")
			f.BoolVar(&cargs.dailySchedule.saturday, "saturday", false, "If set, a backup will be created on Saturdays. Set to false explicitly to clear the flag.")
			f.BoolVar(&cargs.dailySchedule.sunday, "sunday", false, "If set, a backup will be created on Sundays. Set to false explicitly to clear the flag.")
			f.Int32Var(&cargs.timeofday.hours, "hours", 0, "Hours part of the time of day (0-23)")
			f.Int32Var(&cargs.timeofday.minutes, "minutes", 0, "Minutes part of the time of day (0-59)")
			f.StringVar(&cargs.timeofday.timezone, "time-zone", "UTC", "The time-zone this time of day applies to (empty means UTC). Names MUST be exactly as defined in RFC-822.")
			f.Int32Var(&cargs.monthlySchedule.dayOfMonth, "day-of-the-month", 1, "Run the backup on the specified day of the month (1-31)")
			f.StringSliceVar(&cargs.additionalRegionIDs, "additional-region-ids", nil, "Add backup to the specified addition regions")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				id, argsUsed := cmd.OptOption("backup-policy-id", cargs.id, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				backupc := backup.NewBackupServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Select a backup policy to update
				item := selection.MustSelectBackupPolicy(ctx, log, id, backupc)
				// Set changes
				f := c.Flags()
				hasChanges := false
				if f.Changed("name") {
					item.Name = cargs.name
					hasChanges = true
				}
				if f.Changed("description") {
					item.Description = cargs.description
					hasChanges = true
				}
				if f.Changed("upload") {
					item.Upload = cargs.upload
					hasChanges = true
				}
				if f.Changed("paused") {
					item.IsPaused = cargs.paused
					hasChanges = true
				}
				if f.Changed("email-notification") {
					item.EmailNotification = cargs.emailNotification
					hasChanges = true
				}
				if f.Changed("schedule-type") {
					cargs.scheduleType = capitalizeScheduleType(cargs.scheduleType)
					switch cargs.scheduleType {
					case hourly:
						item.Schedule = &backup.BackupPolicy_Schedule{
							HourlySchedule: &backup.BackupPolicy_HourlySchedule{},
							ScheduleType:   hourly,
						}
					case daily:
						item.Schedule = &backup.BackupPolicy_Schedule{
							DailySchedule: &backup.BackupPolicy_DailySchedule{
								ScheduleAt: &backup.TimeOfDay{},
							},
							ScheduleType: daily,
						}
					case monthly:
						item.Schedule = &backup.BackupPolicy_Schedule{
							MonthlySchedule: &backup.BackupPolicy_MonthlySchedule{
								ScheduleAt: &backup.TimeOfDay{},
							},
							ScheduleType: monthly,
						}
					default:
						log.Fatal().Msgf("Invalid schedule type %s", cargs.scheduleType)
					}
					hasChanges = true
				}
				if f.Changed("retention-period") {
					t := time.Duration(cargs.retentionPeriod) * time.Hour
					item.RetentionPeriod = types.DurationProto(t)
					hasChanges = true
				}

				if f.Changed("monday") {
					item.Schedule.DailySchedule.Monday = cargs.dailySchedule.monday
					hasChanges = true
				}
				if f.Changed("tuesday") {
					item.Schedule.DailySchedule.Tuesday = cargs.dailySchedule.tuesday
					hasChanges = true
				}
				if f.Changed("wednesday") {
					item.Schedule.DailySchedule.Wednesday = cargs.dailySchedule.wednesday
					hasChanges = true
				}
				if f.Changed("thursday") {
					item.Schedule.DailySchedule.Thursday = cargs.dailySchedule.thursday
					hasChanges = true
				}
				if f.Changed("friday") {
					item.Schedule.DailySchedule.Friday = cargs.dailySchedule.friday
					hasChanges = true
				}
				if f.Changed("saturday") {
					item.Schedule.DailySchedule.Saturday = cargs.dailySchedule.saturday
					hasChanges = true
				}
				if f.Changed("sunday") {
					item.Schedule.DailySchedule.Sunday = cargs.dailySchedule.sunday
					hasChanges = true
				}
				if f.Changed("locked") {
					item.Locked = cargs.locked
					hasChanges = true
				}
				if f.Changed("hours") {
					switch item.Schedule.ScheduleType {
					case daily:
						item.Schedule.DailySchedule.ScheduleAt.Hours = cargs.timeofday.hours
					case monthly:
						item.Schedule.MonthlySchedule.ScheduleAt.Hours = cargs.timeofday.hours
					}
					hasChanges = true
				}
				if f.Changed("minutes") {
					switch item.Schedule.ScheduleType {
					case daily:
						item.Schedule.DailySchedule.ScheduleAt.Minutes = cargs.timeofday.minutes
					case monthly:
						item.Schedule.MonthlySchedule.ScheduleAt.Minutes = cargs.timeofday.minutes
					}
					hasChanges = true
				}
				if f.Changed("time-zone") {
					switch item.Schedule.ScheduleType {
					case daily:
						item.Schedule.DailySchedule.ScheduleAt.TimeZone = cargs.timeofday.timezone
					case monthly:
						item.Schedule.MonthlySchedule.ScheduleAt.TimeZone = cargs.timeofday.timezone
					}
					hasChanges = true
				}
				if f.Changed("day-of-the-month") {
					item.Schedule.MonthlySchedule.DayOfMonth = cargs.monthlySchedule.dayOfMonth
					hasChanges = true
				}
				if f.Changed("every-interval-hours") {
					item.Schedule.HourlySchedule.ScheduleEveryIntervalHours = cargs.hourlySchedule.scheduleEveryIntervalHours
					hasChanges = true
				}
				if f.Changed("minutes-offset") {
					item.Schedule.HourlySchedule.MinutesOffset = cargs.hourlySchedule.minutesOffset
					hasChanges = true
				}
				if f.Changed("additional-region-ids") {
					item.AdditionalRegionIds = cargs.additionalRegionIDs
					hasChanges = true
				}

				if !hasChanges {
					fmt.Println("No changes")
					return
				}
				// Update backup
				updated, err := backupc.UpdateBackupPolicy(ctx, item)
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to update backup policy")
				}

				// Show result
				fmt.Println("Updated backup policy!")
				fmt.Println(format.BackupPolicy(updated, cmd.RootArgs.Format))
			}
		},
	)
}
