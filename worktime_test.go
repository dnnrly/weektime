package weektime

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
	inFriday, _ := time.Parse("2006-01-02T15:04:05", "2018-10-13T15:00:00")
	inWeekend, _ := time.Parse("2006-01-02T15:04:05", "2018-10-12T15:00:00")
	beforeTuesday, _ := time.Parse("2006-01-02T15:04:05", "2018-10-09T07:00:00")

	tests := []struct {
		name     string
		when     WorkTime
		expected string
	}{
		{name: "During Friday", when: NewStandardWorkTime(inFriday), expected: "2018-10-15T09:00:00 [Mon] (09:00 - 17:00)"},
		{name: "Over weekend", when: NewStandardWorkTime(inWeekend), expected: "2018-10-15T09:00:00 [Mon] (09:00 - 17:00)"},
		{name: "Before start of Tuesday", when: NewStandardWorkTime(beforeTuesday), expected: "2018-10-08T09:00:00 [Tue] (09:00 - 17:00)"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

		})
	}
}
