package util

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// DurationUtil has the
type DurationUtil struct {
	DayFormatExp      *regexp.Regexp
	TimeUnitFormatExp *regexp.Regexp
}

const (
	EmptyDuration time.Duration = time.Duration(0)
)

// ExpandEpochTime "expands" the given time from a string. If it is an int64, it assumes the time
// is a UNIX epoch time and is "absolute" and so refers the time. If it is a string, it assumes
// the time is "relative" to now and returns the UNIX epoch time with the duration added.
func (d *DurationUtil) ExpandEpochTime(str string) (int64, error) {
	// if the incoming time can be used as a number, assume it's an "absolute" time
	epoch, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		return epoch, nil
	}

	// if it's a string, assume that it's "relative" to now, so return that...
	duration, err := d.ParseDuration(str)

	if err != nil {
		return 0, err
	}

	return time.Now().Add(duration).Unix(), nil

}

// ParseDuration accepts a string to parse and return a `time.Duration`.
// This is used because the default time.Duration in go only supports up
// to the hour, and for lease expirations we want to support days,
func (d *DurationUtil) ParseDuration(str string) (time.Duration, error) {

	if d.TimeUnitFormatExp.Match([]byte(str)) {
		// use the default time.Duration behavior
		dur, err := time.ParseDuration(str)
		return dur, err
	}

	// calc for days..
	if d.DayFormatExp.Match([]byte(str)) {
		matches := d.DayFormatExp.FindStringSubmatch(str)
		day, err := strconv.Atoi(matches[1])
		if err != nil {
			return EmptyDuration, err
		}
		if day <= 0 {
			// Well, that's just silly...
			return EmptyDuration, fmt.Errorf("invalid zero or negative date: %d", day)
		}
		dur := time.Duration(day*24) * time.Hour
		return dur, nil
	}

	return EmptyDuration, fmt.Errorf("invalid duration format: %s", str)
}

// NewDurationUtil creates a new `DuractionUtil`
func NewDurationUtil() *DurationUtil {
	durationUtil := &DurationUtil{
		DayFormatExp:      regexp.MustCompile(`([\d]+)(\s*)(?:d)`),
		TimeUnitFormatExp: regexp.MustCompile(`[\d]+(ns|us|ms|s|m|h)`),
	}

	return durationUtil
}
