// Package temporalvalidation provides retained validation adapters.
//
// Deprecated: use github.com/faustbrian/go-temporal/adapters/validation. This
// package remains supported for the longer of 180 days after successor public
// availability and two subsequently published stable root-module minor
// releases.
package temporalvalidation

import (
	adapter "github.com/faustbrian/go-temporal/adapters/validation"
	"github.com/faustbrian/go-temporal/dateperiod"
	"github.com/faustbrian/go-temporal/instant"
	"github.com/faustbrian/go-temporal/timeofday"
	validation "github.com/faustbrian/go-validation"
)

func InstantNonEmpty() validation.Validator[instant.Period] { return adapter.InstantNonEmpty() }
func DateNonEmpty() validation.Validator[dateperiod.Period] { return adapter.DateNonEmpty() }
func DailyNonEmpty() validation.Validator[timeofday.Interval] {
	return adapter.DailyNonEmpty()
}

func TimeBetween(minimum, maximum timeofday.Time) (validation.Validator[timeofday.Time], error) {
	return adapter.TimeBetween(minimum, maximum)
}

func DurationBetween(minimum, maximum timeofday.Duration) (validation.Validator[timeofday.Duration], error) {
	return adapter.DurationBetween(minimum, maximum)
}
