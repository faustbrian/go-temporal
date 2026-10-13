// Package temporalconfig provides retained configuration adapters.
//
// Deprecated: use github.com/faustbrian/go-temporal/v2/adapters/config. This
// package remains supported for the longer of 180 days after successor public
// availability and two subsequently published stable root-module minor
// releases.
package temporalconfig

import (
	temporal "github.com/faustbrian/go-temporal/v2"
	adapter "github.com/faustbrian/go-temporal/v2/adapters/config"
	"github.com/faustbrian/go-temporal/v2/dateperiod"
	"github.com/faustbrian/go-temporal/v2/instant"
	"github.com/faustbrian/go-temporal/v2/timeofday"
)

// InstantPeriod retains the adapters/config instant-period text boundary.
// Its zero value wraps the corresponding core zero value; it is not nullable.
type InstantPeriod struct{ value adapter.InstantPeriod }

// NewInstantPeriod wraps value without parsing or adding validation.
func NewInstantPeriod(value instant.Period) InstantPeriod {
	return InstantPeriod{value: adapter.NewInstantPeriod(value)}
}

// Value returns the wrapped value by value, without exposing mutable storage.
func (v InstantPeriod) Value() instant.Period { return v.value.Value() }

// MarshalText returns owned bytes in ISO 80000 instant notation.
// Formatting is not an independent parse-admission check.
func (v InstantPeriod) MarshalText() ([]byte, error) { return v.value.MarshalText() }

// UnmarshalText parses with default limits and changes the value only on success.
// A nil receiver returns ErrUnsupported; input is not retained. Concurrent reads
// or writes of this wrapper require caller synchronization.
func (v *InstantPeriod) UnmarshalText(text []byte) error {
	if v == nil {
		return temporal.ErrUnsupported
	}
	return v.value.UnmarshalText(text)
}

// DatePeriod retains the adapters/config civil-date period text boundary.
// Its zero value wraps the corresponding core zero value; it is not nullable.
type DatePeriod struct{ value adapter.DatePeriod }

// NewDatePeriod wraps value without parsing or adding validation.
func NewDatePeriod(value dateperiod.Period) DatePeriod {
	return DatePeriod{value: adapter.NewDatePeriod(value)}
}

// Value returns the wrapped value by value, without exposing mutable storage.
func (v DatePeriod) Value() dateperiod.Period { return v.value.Value() }

// MarshalText returns owned bytes in ISO 80000 civil-date notation.
// Formatting is not an independent parse-admission check.
func (v DatePeriod) MarshalText() ([]byte, error) { return v.value.MarshalText() }

// UnmarshalText parses with default limits and changes the value only on success.
// A nil receiver returns ErrUnsupported; input is not retained. Concurrent reads
// or writes of this wrapper require caller synchronization.
func (v *DatePeriod) UnmarshalText(text []byte) error {
	if v == nil {
		return temporal.ErrUnsupported
	}
	return v.value.UnmarshalText(text)
}

// DailyInterval retains the adapters/config daily-interval text boundary.
// Its zero value wraps the corresponding core zero value; it is not nullable.
type DailyInterval struct{ value adapter.DailyInterval }

// NewDailyInterval wraps value without parsing or adding validation.
func NewDailyInterval(value timeofday.Interval) DailyInterval {
	return DailyInterval{value: adapter.NewDailyInterval(value)}
}

// Value returns the wrapped value by value, without exposing mutable storage.
func (v DailyInterval) Value() timeofday.Interval { return v.value.Value() }

// MarshalText returns owned bytes in ISO 80000 daily notation.
// Formatting is not an independent parse-admission check.
func (v DailyInterval) MarshalText() ([]byte, error) { return v.value.MarshalText() }

// UnmarshalText parses with default limits and changes the value only on success.
// A nil receiver returns ErrUnsupported; input is not retained. Concurrent reads
// or writes of this wrapper require caller synchronization.
func (v *DailyInterval) UnmarshalText(text []byte) error {
	if v == nil {
		return temporal.ErrUnsupported
	}
	return v.value.UnmarshalText(text)
}

// Time retains the adapters/config local-time text boundary.
// Its zero value wraps the corresponding core zero value; it is not nullable.
type Time struct{ value adapter.Time }

// NewTime wraps value without parsing or adding validation.
func NewTime(value timeofday.Time) Time { return Time{value: adapter.NewTime(value)} }

// Value returns the wrapped value by value, without exposing mutable storage.
func (v Time) Value() timeofday.Time { return v.value.Value() }

// MarshalText returns owned bytes in the Time.String local-time representation.
// Formatting is not an independent parse-admission check.
func (v Time) MarshalText() ([]byte, error) { return v.value.MarshalText() }

// UnmarshalText parses with default limits and changes the value only on success.
// A nil receiver returns ErrUnsupported; input is not retained. Concurrent reads
// or writes of this wrapper require caller synchronization.
func (v *Time) UnmarshalText(text []byte) error {
	if v == nil {
		return temporal.ErrUnsupported
	}
	return v.value.UnmarshalText(text)
}

// Duration retains the adapters/config fixed-duration text boundary.
// Its zero value wraps the corresponding core zero value; it is not nullable.
type Duration struct{ value adapter.Duration }

// NewDuration wraps value without parsing or adding validation.
func NewDuration(value timeofday.Duration) Duration {
	return Duration{value: adapter.NewDuration(value)}
}

// Value returns the wrapped value by value, without exposing mutable storage.
func (v Duration) Value() timeofday.Duration { return v.value.Value() }

// MarshalText returns owned bytes in the fixed-duration notation.
// Formatting is not an independent parse-admission check.
func (v Duration) MarshalText() ([]byte, error) { return v.value.MarshalText() }

// UnmarshalText parses with default limits and changes the value only on success.
// A nil receiver returns ErrUnsupported; input is not retained. Concurrent reads
// or writes of this wrapper require caller synchronization.
func (v *Duration) UnmarshalText(text []byte) error {
	if v == nil {
		return temporal.ErrUnsupported
	}
	return v.value.UnmarshalText(text)
}
