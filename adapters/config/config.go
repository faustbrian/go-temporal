// Package temporalconfig provides text wrappers for configuration decoders
// that honor encoding.TextUnmarshaler. UnmarshalText assigns only after successful
// parsing; this is transactional assignment, not synchronization for concurrent
// access. It borrows input only during the call and uses default Temporal limits.
// Wrapping a value does not validate it or guarantee that formatting can be parsed
// back. Wrappers own no decoder, goroutine, context, or shutdown lifecycle.
package temporalconfig

import (
	temporal "github.com/faustbrian/go-temporal/v3"
	"github.com/faustbrian/go-temporal/v3/dateperiod"
	"github.com/faustbrian/go-temporal/v3/instant"
	"github.com/faustbrian/go-temporal/v3/internal/diagnostic"
	"github.com/faustbrian/go-temporal/v3/notation"
	"github.com/faustbrian/go-temporal/v3/timeofday"
)

// InstantPeriod is a configuration boundary for an immutable instant period.
type InstantPeriod struct{ value instant.Period }

// NewInstantPeriod wraps value without parsing or adding validation.
func NewInstantPeriod(value instant.Period) InstantPeriod { return InstantPeriod{value: value} }

// Value returns the wrapped value by value, without exposing mutable storage.
func (v InstantPeriod) Value() instant.Period { return v.value }

// MarshalText returns owned bytes in ISO 80000 instant notation.
// Formatting is not an independent parse-admission check.
func (v InstantPeriod) MarshalText() ([]byte, error) {
	encoded, err := notation.FormatInstant(v.value, notation.ISO80000, temporal.Limits{})
	return []byte(encoded), err
}

// UnmarshalText parses with default limits and changes the value only on success.
// A nil receiver returns ErrUnsupported; input is not retained. Concurrent reads
// or writes of this wrapper require caller synchronization.
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

// NewDatePeriod wraps value without parsing or adding validation.
func NewDatePeriod(value dateperiod.Period) DatePeriod { return DatePeriod{value: value} }

// Value returns the wrapped value by value, without exposing mutable storage.
func (v DatePeriod) Value() dateperiod.Period { return v.value }

// MarshalText returns owned bytes in ISO 80000 civil-date notation.
// Formatting is not an independent parse-admission check.
func (v DatePeriod) MarshalText() ([]byte, error) {
	encoded, err := notation.FormatDate(v.value, notation.ISO80000, temporal.Limits{})
	return []byte(encoded), err
}

// UnmarshalText parses with default limits and changes the value only on success.
// A nil receiver returns ErrUnsupported; input is not retained. Concurrent reads
// or writes of this wrapper require caller synchronization.
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

// NewDailyInterval wraps value without parsing or adding validation.
func NewDailyInterval(value timeofday.Interval) DailyInterval { return DailyInterval{value: value} }

// Value returns the wrapped value by value, without exposing mutable storage.
func (v DailyInterval) Value() timeofday.Interval { return v.value }

// MarshalText returns owned bytes in ISO 80000 daily notation.
// Formatting is not an independent parse-admission check.
func (v DailyInterval) MarshalText() ([]byte, error) {
	encoded, err := notation.FormatDailyInterval(v.value, notation.ISO80000, temporal.Limits{})
	return []byte(encoded), err
}

// UnmarshalText parses with default limits and changes the value only on success.
// A nil receiver returns ErrUnsupported; input is not retained. Concurrent reads
// or writes of this wrapper require caller synchronization.
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

// NewTime wraps value without parsing or adding validation.
func NewTime(value timeofday.Time) Time { return Time{value: value} }

// Value returns the wrapped value by value, without exposing mutable storage.
func (v Time) Value() timeofday.Time { return v.value }

// MarshalText returns owned bytes in the Time.String local-time representation.
// Formatting is not an independent parse-admission check.
func (v Time) MarshalText() ([]byte, error) {
	return []byte(v.value.String()), nil
}

// UnmarshalText parses with default limits and changes the value only on success.
// A nil receiver returns ErrUnsupported; input is not retained. Concurrent reads
// or writes of this wrapper require caller synchronization.
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

// NewDuration wraps value without parsing or adding validation.
func NewDuration(value timeofday.Duration) Duration { return Duration{value: value} }

// Value returns the wrapped value by value, without exposing mutable storage.
func (v Duration) Value() timeofday.Duration { return v.value }

// MarshalText returns owned bytes in the fixed-duration notation.
// Formatting is not an independent parse-admission check.
func (v Duration) MarshalText() ([]byte, error) {
	encoded, err := notation.FormatDuration(v.value, temporal.Limits{})
	return []byte(encoded), err
}

// UnmarshalText parses with default limits and changes the value only on success.
// A nil receiver returns ErrUnsupported; input is not retained. Concurrent reads
// or writes of this wrapper require caller synchronization.
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
