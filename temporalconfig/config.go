// Package temporalconfig provides retained configuration adapters.
//
// Deprecated: use github.com/faustbrian/go-temporal/adapters/config. This
// package remains supported for the longer of 180 days after successor public
// availability and two subsequently published stable root-module minor
// releases.
package temporalconfig

import (
	adapter "github.com/faustbrian/go-temporal/adapters/config"
	"github.com/faustbrian/go-temporal/dateperiod"
	"github.com/faustbrian/go-temporal/instant"
	"github.com/faustbrian/go-temporal/timeofday"
)

type InstantPeriod struct{ value adapter.InstantPeriod }

func NewInstantPeriod(value instant.Period) InstantPeriod {
	return InstantPeriod{value: adapter.NewInstantPeriod(value)}
}
func (v InstantPeriod) Value() instant.Period            { return v.value.Value() }
func (v InstantPeriod) MarshalText() ([]byte, error)     { return v.value.MarshalText() }
func (v *InstantPeriod) UnmarshalText(text []byte) error { return v.value.UnmarshalText(text) }

type DatePeriod struct{ value adapter.DatePeriod }

func NewDatePeriod(value dateperiod.Period) DatePeriod {
	return DatePeriod{value: adapter.NewDatePeriod(value)}
}
func (v DatePeriod) Value() dateperiod.Period         { return v.value.Value() }
func (v DatePeriod) MarshalText() ([]byte, error)     { return v.value.MarshalText() }
func (v *DatePeriod) UnmarshalText(text []byte) error { return v.value.UnmarshalText(text) }

type DailyInterval struct{ value adapter.DailyInterval }

func NewDailyInterval(value timeofday.Interval) DailyInterval {
	return DailyInterval{value: adapter.NewDailyInterval(value)}
}
func (v DailyInterval) Value() timeofday.Interval        { return v.value.Value() }
func (v DailyInterval) MarshalText() ([]byte, error)     { return v.value.MarshalText() }
func (v *DailyInterval) UnmarshalText(text []byte) error { return v.value.UnmarshalText(text) }

type Time struct{ value adapter.Time }

func NewTime(value timeofday.Time) Time         { return Time{value: adapter.NewTime(value)} }
func (v Time) Value() timeofday.Time            { return v.value.Value() }
func (v Time) MarshalText() ([]byte, error)     { return v.value.MarshalText() }
func (v *Time) UnmarshalText(text []byte) error { return v.value.UnmarshalText(text) }

type Duration struct{ value adapter.Duration }

func NewDuration(value timeofday.Duration) Duration {
	return Duration{value: adapter.NewDuration(value)}
}
func (v Duration) Value() timeofday.Duration        { return v.value.Value() }
func (v Duration) MarshalText() ([]byte, error)     { return v.value.MarshalText() }
func (v *Duration) UnmarshalText(text []byte) error { return v.value.UnmarshalText(text) }
