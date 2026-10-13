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

// InstantNonEmpty rejects values for which instant.Period.IsEmpty reports true.
// Reports use the supplied Validation context path and limits.
func InstantNonEmpty() validation.Validator[instant.Period] { return adapter.InstantNonEmpty() }

// DateNonEmpty rejects values for which dateperiod.Period.IsEmpty reports true.
// Reports use the supplied Validation context path and limits.
func DateNonEmpty() validation.Validator[dateperiod.Period] { return adapter.DateNonEmpty() }

// DailyNonEmpty rejects CollapsedKind, not every zero-length interval representation.
// Reports use the supplied Validation context path and limits.
func DailyNonEmpty() validation.Validator[timeofday.Interval] {
	return adapter.DailyNonEmpty()
}

// TimeBetween constructs an inclusive, non-circular range validator; reversed bounds fail.
// Reports use the supplied Validation context path and limits.
func TimeBetween(minimum, maximum timeofday.Time) (validation.Validator[timeofday.Time], error) {
	return adapter.TimeBetween(minimum, maximum)
}

// DurationBetween constructs an inclusive elapsed-duration validator; reversed bounds fail.
// Reports use the supplied Validation context path and limits.
func DurationBetween(minimum, maximum timeofday.Duration) (validation.Validator[timeofday.Duration], error) {
	return adapter.DurationBetween(minimum, maximum)
}
