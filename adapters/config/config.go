// Package temporalconfig provides atomic text wrappers for config and other
// configuration decoders that honor encoding.TextUnmarshaler.
package temporalconfig

import (
	temporal "github.com/faustbrian/go-temporal/v2"
	"github.com/faustbrian/go-temporal/v2/dateperiod"
	"github.com/faustbrian/go-temporal/v2/instant"
	"github.com/faustbrian/go-temporal/v2/internal/diagnostic"
	"github.com/faustbrian/go-temporal/v2/notation"
	"github.com/faustbrian/go-temporal/v2/timeofday"
)

// InstantPeriod is a configuration boundary for an immutable instant period.
type InstantPeriod struct{ value instant.Period }

func NewInstantPeriod(value instant.Period) InstantPeriod { return InstantPeriod{value: value} }
func (v InstantPeriod) Value() instant.Period             { return v.value }
func (v InstantPeriod) MarshalText() ([]byte, error) {
	encoded, err := notation.FormatInstant(v.value, notation.ISO80000, temporal.Limits{})
	return []byte(encoded), err
}
func (v *InstantPeriod) UnmarshalText(text []byte) error {
	if v == nil {
		return temporal.ErrUnsupported
	}
	value, err := boundedText(text)
	if err != nil {
		return err
	}
	parsed, err := notation.ParseInstant(value, notation.ISO80000, temporal.Limits{})
	if err != nil {
		return err
	}
	v.value = parsed
	return nil
}

// DatePeriod is a configuration boundary for an immutable civil-date period.
type DatePeriod struct{ value dateperiod.Period }

func NewDatePeriod(value dateperiod.Period) DatePeriod { return DatePeriod{value: value} }
func (v DatePeriod) Value() dateperiod.Period          { return v.value }
func (v DatePeriod) MarshalText() ([]byte, error) {
	encoded, err := notation.FormatDate(v.value, notation.ISO80000, temporal.Limits{})
	return []byte(encoded), err
}
func (v *DatePeriod) UnmarshalText(text []byte) error {
	if v == nil {
		return temporal.ErrUnsupported
	}
	value, err := boundedText(text)
	if err != nil {
		return err
	}
	parsed, err := notation.ParseDate(value, notation.ISO80000, temporal.Limits{})
	if err != nil {
		return err
	}
	v.value = parsed
	return nil
}

// DailyInterval is a configuration boundary for a daily interval.
type DailyInterval struct{ value timeofday.Interval }

func NewDailyInterval(value timeofday.Interval) DailyInterval { return DailyInterval{value: value} }
func (v DailyInterval) Value() timeofday.Interval             { return v.value }
func (v DailyInterval) MarshalText() ([]byte, error) {
	encoded, err := notation.FormatDailyInterval(v.value, notation.ISO80000, temporal.Limits{})
	return []byte(encoded), err
}
func (v *DailyInterval) UnmarshalText(text []byte) error {
	if v == nil {
		return temporal.ErrUnsupported
	}
	value, err := boundedText(text)
	if err != nil {
		return err
	}
	parsed, err := notation.ParseDailyInterval(value, notation.ISO80000, temporal.Limits{})
	if err != nil {
		return err
	}
	v.value = parsed
	return nil
}

// Time is a configuration boundary for a local time-of-day value.
type Time struct{ value timeofday.Time }

func NewTime(value timeofday.Time) Time { return Time{value: value} }
func (v Time) Value() timeofday.Time    { return v.value }
func (v Time) MarshalText() ([]byte, error) {
	return []byte(v.value.String()), nil
}
func (v *Time) UnmarshalText(text []byte) error {
	if v == nil {
		return temporal.ErrUnsupported
	}
	value, err := boundedText(text)
	if err != nil {
		return err
	}
	parsed, err := timeofday.Parse(value, temporal.Limits{})
	if err != nil {
		return err
	}
	v.value = parsed
	return nil
}

// Duration is a configuration boundary for a fixed elapsed duration.
type Duration struct{ value timeofday.Duration }

func NewDuration(value timeofday.Duration) Duration { return Duration{value: value} }
func (v Duration) Value() timeofday.Duration        { return v.value }
func (v Duration) MarshalText() ([]byte, error) {
	encoded, err := notation.FormatDuration(v.value, temporal.Limits{})
	return []byte(encoded), err
}
func (v *Duration) UnmarshalText(text []byte) error {
	if v == nil {
		return temporal.ErrUnsupported
	}
	value, err := boundedText(text)
	if err != nil {
		return err
	}
	parsed, err := notation.ParseDuration(value, temporal.Limits{})
	if err != nil {
		return err
	}
	v.value = parsed
	return nil
}

func boundedText(text []byte) (string, error) {
	limits := temporal.DefaultLimits()
	if len(text) > limits.ParseBytes {
		cause := &temporal.LimitError{Field: "parse_bytes", Value: len(text), Max: limits.ParseBytes}
		return "", diagnostic.New(limits.ErrorBytes, temporal.ErrLimit.Error(), temporal.ErrLimit, cause)
	}
	return string(text), nil
}
