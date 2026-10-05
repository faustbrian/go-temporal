// Package temporalvalidation provides retained validation adapters.
//
// Deprecated: use github.com/faustbrian/go-temporal/v2/adapters/validation. This
// package remains supported for the longer of 180 days after successor public
// availability and two subsequently published stable root-module minor
// releases.
package temporalvalidation

import (
	adapter "github.com/faustbrian/go-temporal/v2/adapters/validation"
	"github.com/faustbrian/go-temporal/v2/dateperiod"
	"github.com/faustbrian/go-temporal/v2/instant"
	"github.com/faustbrian/go-temporal/v2/timeofday"
	validation "github.com/faustbrian/go-validation/v2"
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
