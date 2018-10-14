package worktime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkTime_String(t *testing.T) {
	weekend, _ := time.Parse("2006-01-02T15:04:05", "2018-10-13T15:00:00")
	inHours, _ := time.Parse("2006-01-02T15:04:05", "2018-10-12T15:00:00")
	outHours, _ := time.Parse("2006-01-02T15:04:05", "2018-10-12T20:00:00")

	tests := []struct {
		name     string
		when     WorkTime
		expected string
	}{
		{name: "In work hours", when: NewStandardWorkTime(inHours), expected: "2018-10-12T15:00:00 [Fri] (09:00 - 17:00)"},
		{name: "On weekend", when: NewStandardWorkTime(weekend), expected: "2018-10-13T15:00:00 [Sat] (09:00 - 17:00)"},
		{name: "Work day, out of hours", when: NewStandardWorkTime(outHours), expected: "2018-10-12T20:00:00 [Fri] (09:00 - 17:00)"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.when.String())
		})
	}
}

func TestWorkTime_NextStart(t *testing.T) {
	inFriday, _ := time.Parse("2006-01-02T15:04:05", "2018-10-12T15:00:00")
	inSaturday, _ := time.Parse("2006-01-02T15:04:05", "2018-10-13T15:00:00")
	inSunday, _ := time.Parse("2006-01-02T15:04:05", "2018-10-14T15:00:00")
	beforeTuesday, _ := time.Parse("2006-01-02T15:04:05", "2018-10-09T07:00:00")

	tests := []struct {
		name     string
		when     WorkTime
		expected string
	}{
		{name: "During Friday", when: NewStandardWorkTime(inFriday), expected: "2018-10-15T09:00:00 [Mon] (09:00 - 17:00)"},
		{name: "During Saturday", when: NewStandardWorkTime(inSaturday), expected: "2018-10-15T09:00:00 [Mon] (09:00 - 17:00)"},
		{name: "During Sunday", when: NewStandardWorkTime(inSunday), expected: "2018-10-15T09:00:00 [Mon] (09:00 - 17:00)"},
		{name: "Before start of Tuesday", when: NewStandardWorkTime(beforeTuesday), expected: "2018-10-09T09:00:00 [Tue] (09:00 - 17:00)"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.when.NextStart().String())
		})
	}
}

func TestNewStandardWorkTime(t *testing.T) {
	wt := NewStandardWorkTime(time.Now())
	assert.Equal(t, (time.Hour * 9).String(), wt.start.String())
	assert.Equal(t, (time.Hour * 17).String(), wt.end.String())
}

func TestNewWorkTime(t *testing.T) {
	wt := NewWorkTime(time.Now(), time.Hour*8, time.Hour*18)
	assert.Equal(t, (time.Hour * 8).String(), wt.start.String())
	assert.Equal(t, (time.Hour * 18).String(), wt.end.String())
}

func TestWorkTime_IsWorkDay(t *testing.T) {
	sun, _ := time.Parse("2006-01-02T15:04:05", "2018-10-07T15:00:00")
	mon, _ := time.Parse("2006-01-02T15:04:05", "2018-10-08T15:00:00")
	tue, _ := time.Parse("2006-01-02T15:04:05", "2018-10-09T15:00:00")
	wed, _ := time.Parse("2006-01-02T15:04:05", "2018-10-10T15:00:00")
	thu, _ := time.Parse("2006-01-02T15:04:05", "2018-10-11T15:00:00")
	fri, _ := time.Parse("2006-01-02T15:04:05", "2018-10-12T15:00:00")
	sat, _ := time.Parse("2006-01-02T15:04:05", "2018-10-13T15:00:00")

	tests := []struct {
		name string
		wt   WorkTime
		want bool
	}{
		{name: "Sun", wt: NewStandardWorkTime(sun), want: false},
		{name: "Mon", wt: NewStandardWorkTime(mon), want: true},
		{name: "Tue", wt: NewStandardWorkTime(tue), want: true},
		{name: "Wed", wt: NewStandardWorkTime(wed), want: true},
		{name: "Thu", wt: NewStandardWorkTime(thu), want: true},
		{name: "Fri", wt: NewStandardWorkTime(fri), want: true},
		{name: "Sat", wt: NewStandardWorkTime(sat), want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.wt.IsWorkDay(); got != tt.want {
				t.Errorf("WorkTime.IsWorkDay() = %v, want %v", got, tt.want)
			}
			if day := tt.wt.Format("Mon"); day != tt.name {
				t.Errorf("WorkTime Day = %v, want %v", day, tt.name)
			}
		})
	}
}

func TestWorkTime_Length(t *testing.T) {
	tests := []struct {
		name  string
		start time.Duration
		end   time.Duration
		want  time.Duration
	}{
		{name: "Normal", start: time.Hour * 9, end: time.Hour * 17, want: time.Hour * 8},
		{name: "Early shift", start: time.Hour * 6, end: time.Hour * 12, want: time.Hour * 6},
		{name: "Late shift", start: time.Hour * 19, end: time.Hour * 03, want: time.Hour * 8},
		{name: "Day", start: time.Hour * 0, end: time.Hour * 0, want: time.Hour * 24},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wt := WorkTime{
				Time:  time.Now(),
				start: tt.start,
				end:   tt.end,
			}
			if got := wt.Length(); got != tt.want {
				t.Errorf("WorkTime.Length() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkTime_SinceMidnight(t *testing.T) {
	when, _ := time.Parse("2006-01-02T15:04:05", "2018-10-07T15:00:00")
	wt := NewStandardWorkTime(when)
	assert.Equal(t, time.Hour*15, wt.SinceMidnight())
}

func TestWorkTime_BeforeStart(t *testing.T) {
	before, _ := time.Parse("2006-01-02T15:04:05", "2018-10-08T07:00:00")
	during, _ := time.Parse("2006-01-02T15:04:05", "2018-10-08T15:00:00")
	after, _ := time.Parse("2006-01-02T15:04:05", "2018-10-08T19:00:00")

	type fields struct {
		Time  time.Time
		start time.Duration
		end   time.Duration
	}
	tests := []struct {
		name string
		Time time.Time
		want bool
	}{
		{name: "Before", Time: before, want: true},
		{name: "During", Time: during, want: false},
		{name: "After", Time: after, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wt := NewStandardWorkTime(tt.Time)
			if got := wt.BeforeStart(); got != tt.want {
				t.Errorf("WorkTime.BeforeStart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkTime_DuringOfficeHours(t *testing.T) {
	before, _ := time.Parse("2006-01-02T15:04:05", "2018-10-12T07:00:00")
	weekend, _ := time.Parse("2006-01-02T15:04:05", "2018-10-13T15:00:00")
	during, _ := time.Parse("2006-01-02T15:04:05", "2018-10-12T15:00:00")
	after, _ := time.Parse("2006-01-02T15:04:05", "2018-10-12T19:00:00")

	tests := []struct {
		name string
		Time time.Time
		want bool
	}{
		{name: "Before", Time: before, want: false},
		{name: "During", Time: during, want: true},
		{name: "After", Time: after, want: false},
		{name: "Weekend", Time: weekend, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wt := NewStandardWorkTime(tt.Time)
			if got := wt.DuringOfficeHours(); got != tt.want {
				t.Errorf("WorkTime.DuringOfficeHours() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkTime_AfterEnd(t *testing.T) {
	before, _ := time.Parse("2006-01-02T15:04:05", "2018-10-08T07:00:00")
	during, _ := time.Parse("2006-01-02T15:04:05", "2018-10-08T15:00:00")
	after, _ := time.Parse("2006-01-02T15:04:05", "2018-10-08T19:00:00")

	type fields struct {
		Time  time.Time
		start time.Duration
		end   time.Duration
	}
	tests := []struct {
		name string
		Time time.Time
		want bool
	}{
		{name: "Before", Time: before, want: false},
		{name: "During", Time: during, want: false},
		{name: "After", Time: after, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wt := NewStandardWorkTime(tt.Time)
			if got := wt.AfterEnd(); got != tt.want {
				t.Errorf("WorkTime.AfterEnd() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestWorkTime_Start(t *testing.T) {
	before, _ := time.Parse("2006-01-02T15:04:05", "2018-10-11T07:00:00")
	during, _ := time.Parse("2006-01-02T15:04:05", "2018-10-11T15:00:00")
	after, _ := time.Parse("2006-01-02T15:04:05", "2018-10-11T19:00:00")
	friday, _ := time.Parse("2006-01-02T15:04:05", "2018-10-12T15:00:00")
	weekend, _ := time.Parse("2006-01-02T15:04:05", "2018-10-13T15:00:00")

	s1, _ := time.Parse("2006-01-02T15:04:05", "2018-10-11T09:00:00")
	s2, _ := time.Parse("2006-01-02T15:04:05", "2018-10-12T09:00:00")
	s3, _ := time.Parse("2006-01-02T15:04:05", "2018-10-13T09:00:00")

	tests := []struct {
		name string
		Time time.Time
		want time.Time
	}{
		{name: "Early", Time: before, want: s1},
		{name: "During", Time: during, want: s1},
		{name: "After", Time: after, want: s1},
		{name: "Friday", Time: friday, want: s2},
		{name: "Weekend", Time: weekend, want: s3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wt := NewStandardWorkTime(tt.Time)
			if got := wt.Start(); !got.Time.Equal(tt.want) {
				t.Errorf("WorkTime.Start() = %v, want %v", got.Time, tt.want)
			}
		})
	}
}

func TestWorkTime_End(t *testing.T) {
	before, _ := time.Parse("2006-01-02T15:04:05", "2018-10-11T07:00:00")
	during, _ := time.Parse("2006-01-02T15:04:05", "2018-10-11T15:00:00")
	after, _ := time.Parse("2006-01-02T15:04:05", "2018-10-11T19:00:00")
	friday, _ := time.Parse("2006-01-02T15:04:05", "2018-10-12T18:00:00")
	weekend, _ := time.Parse("2006-01-02T15:04:05", "2018-10-13T15:00:00")

	s1, _ := time.Parse("2006-01-02T15:04:05", "2018-10-11T17:00:00")
	s2, _ := time.Parse("2006-01-02T15:04:05", "2018-10-12T17:00:00")
	s3, _ := time.Parse("2006-01-02T15:04:05", "2018-10-13T17:00:00")

	tests := []struct {
		name string
		Time time.Time
		want time.Time
	}{
		{name: "Early", Time: before, want: s1},
		{name: "During", Time: during, want: s1},
		{name: "After", Time: after, want: s1},
		{name: "Friday", Time: friday, want: s2},
		{name: "Weekend", Time: weekend, want: s3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wt := NewStandardWorkTime(tt.Time)
			if got := wt.End(); !got.Time.Equal(tt.want) {
				t.Errorf("WorkTime.End() = %v, want %v", got.Time, tt.want)
			}
		})
	}
}

func TestWorkTime_FromStart(t *testing.T) {
	before, _ := time.Parse("2006-01-02T15:04:05", "2018-10-11T07:00:00")
	during, _ := time.Parse("2006-01-02T15:04:05", "2018-10-11T15:00:00")
	after, _ := time.Parse("2006-01-02T15:04:05", "2018-10-11T19:00:00")
	friday, _ := time.Parse("2006-01-02T15:04:05", "2018-10-12T18:00:00")
	weekend, _ := time.Parse("2006-01-02T15:04:05", "2018-10-13T15:00:00")
	tests := []struct {
		name string
		Time time.Time
		want time.Duration
	}{
		{name: "Early", Time: before, want: time.Hour * -2},
		{name: "During", Time: during, want: time.Hour * 6},
		{name: "After", Time: after, want: time.Hour * 10},
		{name: "Friday", Time: friday, want: time.Hour * 9},
		{name: "Weekend", Time: weekend, want: time.Hour * 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wt := NewStandardWorkTime(tt.Time)
			if got := wt.FromStart(); got != tt.want {
				t.Errorf("WorkTime.FromStart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkTime_UntilEnd(t *testing.T) {
	before, _ := time.Parse("2006-01-02T15:04:05", "2018-10-11T07:00:00")
	during, _ := time.Parse("2006-01-02T15:04:05", "2018-10-11T15:00:00")
	after, _ := time.Parse("2006-01-02T15:04:05", "2018-10-11T19:00:00")
	friday, _ := time.Parse("2006-01-02T15:04:05", "2018-10-12T18:00:00")
	weekend, _ := time.Parse("2006-01-02T15:04:05", "2018-10-13T15:00:00")
	tests := []struct {
		name string
		Time time.Time
		want time.Duration
	}{
		{name: "Early", Time: before, want: time.Hour * 10},
		{name: "During", Time: during, want: time.Hour * 2},
		{name: "After", Time: after, want: time.Hour * -2},
		{name: "Friday", Time: friday, want: time.Hour * -1},
		{name: "Weekend", Time: weekend, want: time.Hour * 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wt := NewStandardWorkTime(tt.Time)
			if got := wt.UntilEnd(); got != tt.want {
				t.Errorf("WorkTime.UntilEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkTime_Add(t *testing.T) {
	when, _ := time.Parse("2006-01-02T15:04:05", "2018-10-12T10:00:00")
	expected, _ := time.Parse("2006-01-02T15:04:05", "2018-10-12T21:00:00")

	wt := NewWorkTime(when, time.Hour*7, time.Hour*14)

	result := wt.Add(time.Hour * 11)
	assert.Equal(t, expected, result.Time)
	assert.Equal(t, time.Hour*7, result.start)
	assert.Equal(t, time.Hour*14, result.end)
}

func TestWorkTime_IsSameDay(t *testing.T) {
	before, _ := time.Parse("2006-01-02T15:04:05", "2018-10-11T07:00:00")
	during, _ := time.Parse("2006-01-02T15:04:05", "2018-10-11T15:00:00")
	after, _ := time.Parse("2006-01-02T15:04:05", "2018-10-11T19:00:00")
	friday, _ := time.Parse("2006-01-02T15:04:05", "2018-10-12T18:00:00")
	weekend, _ := time.Parse("2006-01-02T15:04:05", "2018-10-13T15:00:00")

	assert.True(t, NewStandardWorkTime(before).IsSameDay(during))
	assert.True(t, NewStandardWorkTime(before).IsSameDay(after))
	assert.False(t, NewStandardWorkTime(friday).IsSameDay(weekend))
}
