// Package postgres provides retained PostgreSQL adapters.
//
// Deprecated: use github.com/faustbrian/go-temporal/v2/adapters/postgres. This
// package remains supported for the longer of 180 days after successor public
// availability and two subsequently published stable root-module minor
// releases.
//
// Conversions delegate to the canonical adapter without owning a connection or
// transaction. Borrowed inputs and nullable scanner assignment follow its limits;
// successful value conversion is not live-server interoperability qualification.
package postgres

import (
	"database/sql/driver"
	"time"

	calendar "github.com/faustbrian/go-calendar/v2"
	temporal "github.com/faustbrian/go-temporal/v2"
	temporalpostgres "github.com/faustbrian/go-temporal/v2/adapters/postgres"
	"github.com/faustbrian/go-temporal/v2/dateperiod"
	"github.com/faustbrian/go-temporal/v2/instant"
	"github.com/jackc/pgx/v5/pgtype"
)

// InstantRangeValue converts a finite nonempty period without rounding sub-microsecond endpoints.
func InstantRangeValue(period instant.Period) (pgtype.Range[time.Time], error) {
	return temporalpostgres.InstantRangeValue(period)
}

// InstantPeriod rejects NULL, infinite, empty, invalid, and sub-microsecond pgx ranges.
func InstantPeriod(value pgtype.Range[time.Time]) (instant.Period, error) {
	return temporalpostgres.InstantPeriod(value)
}

// InstantMultirangeValue returns an owned pgx range slice, or nil on an element conversion error.
func InstantMultirangeValue(set instant.Set) (pgtype.Multirange[pgtype.Range[time.Time]], error) {
	return temporalpostgres.InstantMultirangeValue(set)
}

// InstantSet admits the borrowed input count before copying and normalizing under limits.
func InstantSet(value pgtype.Multirange[pgtype.Range[time.Time]], limits temporal.Limits) (instant.Set, error) {
	return temporalpostgres.InstantSet(value, limits)
}

// DateRangeValue canonicalizes represented dates to [first, successor(last)); an unavailable successor fails.
func DateRangeValue(period dateperiod.Period) (pgtype.Range[calendar.Date], error) {
	return temporalpostgres.DateRangeValue(period)
}

// DatePeriod normalizes finite nonempty pgx bounds to represented closed civil dates.
func DatePeriod(value pgtype.Range[calendar.Date]) (dateperiod.Period, error) {
	return temporalpostgres.DatePeriod(value)
}

// DateMultirangeValue returns owned canonical pgx date ranges, or nil on conversion failure.
func DateMultirangeValue(set dateperiod.Set) (pgtype.Multirange[pgtype.Range[calendar.Date]], error) {
	return temporalpostgres.DateMultirangeValue(set)
}

// DateSet admits the borrowed input count before copying and normalizing under limits.
func DateSet(value pgtype.Multirange[pgtype.Range[calendar.Date]], limits temporal.Limits) (dateperiod.Set, error) {
	return temporalpostgres.DateSet(value, limits)
}

// InstantRange is a retained nullable SQL wrapper; the zero value represents NULL.
// Mutating Scan calls need caller synchronization. No database resource is owned.
type InstantRange struct{ value temporalpostgres.InstantRange }

// NewInstantRange validates representability and returns a non-NULL SQL wrapper on success.
func NewInstantRange(period instant.Period) (InstantRange, error) {
	value, err := temporalpostgres.NewInstantRange(period)
	return InstantRange{value: value}, err
}

// Period returns the core value and whether this wrapper is non-NULL.
func (r InstantRange) Period() (instant.Period, bool) { return r.value.Period() }

// Value returns SQL NULL for the zero wrapper, otherwise canonical range text.
func (r InstantRange) Value() (driver.Value, error) { return r.value.Value() }

// Scan accepts NULL, string, or bytes with the canonical adapter byte ceiling.
// NULL clears the wrapper; other failures preserve it. A nil receiver fails.
// Borrowed byte input is not retained after the call.
func (r *InstantRange) Scan(source any) error {
	if r == nil {
		return temporal.ErrUnsupported
	}
	return r.value.Scan(source)
}

// DateRange is a retained nullable SQL wrapper; the zero value represents NULL.
// Mutating Scan calls need caller synchronization. No database resource is owned.
type DateRange struct{ value temporalpostgres.DateRange }

// NewDateRange validates and canonicalizes dates before returning a non-NULL SQL wrapper.
func NewDateRange(period dateperiod.Period) (DateRange, error) {
	value, err := temporalpostgres.NewDateRange(period)
	return DateRange{value: value}, err
}

// Period returns the core value and whether this wrapper is non-NULL.
func (r DateRange) Period() (dateperiod.Period, bool) { return r.value.Period() }

// Value returns SQL NULL for the zero wrapper, otherwise canonical range text.
func (r DateRange) Value() (driver.Value, error) { return r.value.Value() }

// Scan accepts NULL, string, or bytes with the canonical adapter byte ceiling.
// NULL clears the wrapper; other failures preserve it. A nil receiver fails.
// Borrowed byte input is not retained after the call.
func (r *DateRange) Scan(source any) error {
	if r == nil {
		return temporal.ErrUnsupported
	}
	return r.value.Scan(source)
}
