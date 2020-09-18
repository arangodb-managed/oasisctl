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

package data

import (
	"fmt"
	"time"
	"unicode"

	"github.com/gogo/protobuf/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	backup "github.com/arangodb-managed/apis/backup/v1"
	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
)

const (
	// Schedule Types
	hourly  = "Hourly"
	daily   = "Daily"
	monthly = "Monthly"
)

func init() {
	cmd.InitCommand(
		createBackupCmd,
		&cobra.Command{
			Use:   "policy",
			Short: "Create a new backup policy",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				name              string
				deploymentID      string
				description       string
				emailNotification string
				scheduleType      string
				paused            bool
				hourlySchedule    struct {
					scheduleEveryIntervalHours int32
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
				retentionPeriod int
				upload          bool
				locked          bool
			}{}
			f.StringVar(&cargs.name, "name", "", "Name of the deployment")
			f.StringVar(&cargs.deploymentID, "deployment-id", "", "ID of the deployment")
			f.StringVar(&cargs.description, "description", "", "Description of the backup policy")
			f.StringVar(&cargs.emailNotification, "email-notification", "", "Email notification setting (Never|FailureOnly|Always)")
			f.StringVar(&cargs.scheduleType, "schedule-type", "", "Schedule of the policy (Hourly|Daily|Monthly)")
			f.BoolVar(&cargs.upload, "upload", false, "The backup should be uploaded")
			f.BoolVar(&cargs.paused, "paused", false, "The policy is paused")
			f.IntVar(&cargs.retentionPeriod, "retention-period", 0, "Backups created by this policy will be automatically deleted after the specified retention period. A value of 0 means that backup will never be deleted.")
			f.Int32Var(&cargs.hourlySchedule.scheduleEveryIntervalHours, "every-interval-hours", 0, "Schedule should run with an interval of the specified hours (1-23)")
			f.BoolVar(&cargs.dailySchedule.monday, "monday", false, "If set, a backup will be created on Mondays")
			f.BoolVar(&cargs.dailySchedule.tuesday, "tuesday", false, "If set, a backup will be created on Tuesdays")
			f.BoolVar(&cargs.dailySchedule.wednesday, "wednesday", false, "If set, a backup will be created on Wednesdays")
			f.BoolVar(&cargs.dailySchedule.thursday, "thursday", false, "If set, a backup will be created on Thursdays")
			f.BoolVar(&cargs.dailySchedule.friday, "friday", false, "If set, a backup will be created on Fridays.")
			f.BoolVar(&cargs.dailySchedule.saturday, "saturday", false, "If set, a backup will be created on Saturdays")
			f.BoolVar(&cargs.dailySchedule.sunday, "sunday", false, "If set, a backup will be created on Sundays")
			f.Int32Var(&cargs.timeofday.hours, "hours", 0, "Hours part of the time of day (0-23)")
			f.Int32Var(&cargs.timeofday.minutes, "minutes", 0, "Minutes part of the time of day (0-59)")
			f.StringVar(&cargs.timeofday.timezone, "time-zone", "UTC", "The time-zone this time of day applies to (empty means UTC). Names MUST be exactly as defined in RFC-822.")
			f.Int32Var(&cargs.monthlySchedule.dayOfMonth, "day-of-the-month", 1, "Run the backup on the specified day of the month (1-31)")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				name, argsUsed := cmd.ReqOption("name", cargs.name, args, 0)
				deploymentID, argsUsed := cmd.ReqOption("deployment-id", cargs.deploymentID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				backupc := backup.NewBackupServiceClient(conn)
				ctx := cmd.ContextWithToken()

				cargs.scheduleType = capitalizeScheduleType(cargs.scheduleType)
				b := &backup.BackupPolicy{
					Name:         name,
					Description:  cargs.description,
					DeploymentId: deploymentID,
					Schedule: &backup.BackupPolicy_Schedule{
						ScheduleType: cargs.scheduleType,
					},
					Upload:            cargs.upload,
					EmailNotification: cargs.emailNotification,
					IsPaused:          cargs.paused,
				}

				switch cargs.scheduleType {
				case hourly:
					b.Schedule.HourlySchedule = &backup.BackupPolicy_HourlySchedule{
						ScheduleEveryIntervalHours: cargs.hourlySchedule.scheduleEveryIntervalHours,
					}
				case daily:
					b.Schedule.DailySchedule = &backup.BackupPolicy_DailySchedule{
						Monday:    cargs.dailySchedule.monday,
						Tuesday:   cargs.dailySchedule.tuesday,
						Wednesday: cargs.dailySchedule.wednesday,
						Thursday:  cargs.dailySchedule.thursday,
						Friday:    cargs.dailySchedule.friday,
						Saturday:  cargs.dailySchedule.saturday,
						Sunday:    cargs.dailySchedule.sunday,
						ScheduleAt: &backup.TimeOfDay{
							Hours:    cargs.timeofday.hours,
							Minutes:  cargs.timeofday.minutes,
							TimeZone: cargs.timeofday.timezone,
						},
					}
				case monthly:
					b.Schedule.MonthlySchedule = &backup.BackupPolicy_MonthlySchedule{
						DayOfMonth: cargs.monthlySchedule.dayOfMonth,
						ScheduleAt: &backup.TimeOfDay{
							Hours:    cargs.timeofday.hours,
							Minutes:  cargs.timeofday.minutes,
							TimeZone: cargs.timeofday.timezone,
						},
					}
				default:
					log.Fatal().Msgf("Invalid schedule type %s", cargs.scheduleType)
				}

				t := time.Duration(cargs.retentionPeriod) * time.Hour
				b.RetentionPeriod = types.DurationProto(t)

				result, err := backupc.CreateBackupPolicy(ctx, b)

				if err != nil {
					log.Fatal().Err(err).Msg("Failed to create backup policy")
				}

				// Show result
				format.DisplaySuccess(cmd.RootArgs.Format)
				fmt.Println(format.BackupPolicy(result, cmd.RootArgs.Format))
			}
		},
	)
}

// capitalizeScheduleType creates Daily, Hourly, Monthly out of the uncapitalized words.
func capitalizeScheduleType(t string) string {
	head := t[0]
	tail := t[1:]
	h := unicode.ToUpper(rune(head))
	return string(h) + tail
}
