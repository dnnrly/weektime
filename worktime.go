// Package worktime extends time.Time to help with work times and working week calculations
package worktime

import (
	"fmt"
	"time"
)

// Day is a full calendar day
const Day = time.Hour * 24

// WorkTime represents a day that could be a normal week day with office hours
type WorkTime struct {
	time.Time

	start time.Duration
	end   time.Duration
}

// NewStandardWorkTime returns a WorkTime with 9-5 as the working hours
func NewStandardWorkTime(t time.Time) WorkTime {
	return NewWorkTime(
		t,
		time.Hour*9,
		time.Hour*17,
	)
}

// NewWorkTime allows you to specify the start and end of the work day
func NewWorkTime(t time.Time, start, end time.Duration) WorkTime {
	return WorkTime{
		Time:  t,
		start: start,
		end:   end,
	}
}

func (t WorkTime) String() string {
	return fmt.Sprintf("%s (%02.0f:%02.0f - %02.0f:%02.0f)",
		t.Format("2006-01-02T15:04:05 [Mon]"),
		t.start.Hours(), t.start.Minutes()-t.start.Hours()*60,
		t.end.Hours(), t.end.Minutes()-t.end.Hours()*60,
	)
}

// IsWorkDay returns false if the time is a Saturday or Sunday, else true
func (t WorkTime) IsWorkDay() bool {
	day := t.Weekday()
	if day == time.Sunday || day == time.Saturday {
		return false
	}

	return true
}

// Length returns the time.Duration of the working day
func (t WorkTime) Length() time.Duration {
	diff := t.end - t.start

	if diff <= 0 {
		diff += Day
	}

	return diff
}

// SinceMidnight returns the time.Duration since the previous midnight
func (t WorkTime) SinceMidnight() time.Duration {
	return t.Sub(t.Truncate(Day))
}

// BeforeStart returns true if the time of day is before the start of the working day
func (t WorkTime) BeforeStart() bool {
	return t.SinceMidnight() < t.start
}

// DuringOfficeHours returns true if the time is during the working hours
func (t WorkTime) DuringOfficeHours() bool {
	return t.IsWorkDay() && !t.BeforeStart() && !t.AfterEnd()
}

// AfterEnd returns true if the time of day is after the end of the working day
func (t WorkTime) AfterEnd() bool {
	return t.SinceMidnight() > t.end
}

// Start returns a WorkTime representing the working start of the current day
func (t WorkTime) Start() WorkTime {
	return NewWorkTime(
		t.Truncate(Day).Add(t.start),
		t.start,
		t.end,
	)
}

// NextStart returns the start of the next working day
func (t WorkTime) NextStart() WorkTime {
	if !t.BeforeStart() {
		d := Day
		if t.Weekday() == time.Friday {
			d = Day * 3
		} else if t.Weekday() == time.Saturday {
			d = Day * 2
		}

		t = t.Add(d)
	}

	return t.Start()
}

// End returns a WorkTime representing the end of the working day
func (t WorkTime) End() WorkTime {
	return NewWorkTime(
		t.Truncate(Day).Add(t.end),
		t.start,
		t.end,
	)
}

// FromStart returns the amount of time since the start of the working day, can be negative if the
// time of day is before the start of the working day
func (t WorkTime) FromStart() time.Duration {
	return t.SinceMidnight() - t.start
}

// UntilEnd returns the amount of time until the end of the working day, can be negative if the
// time of day is after the end of the working day
func (t WorkTime) UntilEnd() time.Duration {
	return t.SinceMidnight() - t.end
}

// Add returns a WorkTime plus a time.Duration
func (t WorkTime) Add(d time.Duration) WorkTime {
	return NewWorkTime(t.Time.Add(d), t.start, t.end)
}

// IsSameDay returns true if this WorkTime is the same day as d
func (t WorkTime) IsSameDay(d time.Time) bool {
	tYear, tMonth, tDay := t.Date()
	dYear, dMonth, dDay := d.Date()

	return tYear == dYear &&
		tMonth == dMonth &&
		tDay == dDay
}
