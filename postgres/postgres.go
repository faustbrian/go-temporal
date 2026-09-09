// Package postgres provides retained PostgreSQL adapters.
//
// Deprecated: use github.com/faustbrian/go-temporal/adapters/postgres. This
// package remains supported for the longer of 180 days after successor public
// availability and two subsequently published stable root-module minor
// releases.
package postgres

import (
	"database/sql/driver"
	"time"

	calendar "github.com/faustbrian/go-calendar"
	temporal "github.com/faustbrian/go-temporal"
	temporalpostgres "github.com/faustbrian/go-temporal/adapters/postgres"
	"github.com/faustbrian/go-temporal/dateperiod"
	"github.com/faustbrian/go-temporal/instant"
	"github.com/jackc/pgx/v5/pgtype"
)

func InstantRangeValue(period instant.Period) (pgtype.Range[time.Time], error) {
	return temporalpostgres.InstantRangeValue(period)
}
func InstantPeriod(value pgtype.Range[time.Time]) (instant.Period, error) {
	return temporalpostgres.InstantPeriod(value)
}
func InstantMultirangeValue(set instant.Set) (pgtype.Multirange[pgtype.Range[time.Time]], error) {
	return temporalpostgres.InstantMultirangeValue(set)
}
func InstantSet(value pgtype.Multirange[pgtype.Range[time.Time]], limits temporal.Limits) (instant.Set, error) {
	return temporalpostgres.InstantSet(value, limits)
}
func DateRangeValue(period dateperiod.Period) (pgtype.Range[calendar.Date], error) {
	return temporalpostgres.DateRangeValue(period)
}
func DatePeriod(value pgtype.Range[calendar.Date]) (dateperiod.Period, error) {
	return temporalpostgres.DatePeriod(value)
}
func DateMultirangeValue(set dateperiod.Set) (pgtype.Multirange[pgtype.Range[calendar.Date]], error) {
	return temporalpostgres.DateMultirangeValue(set)
}
func DateSet(value pgtype.Multirange[pgtype.Range[calendar.Date]], limits temporal.Limits) (dateperiod.Set, error) {
	return temporalpostgres.DateSet(value, limits)
}

type InstantRange struct{ value temporalpostgres.InstantRange }

func NewInstantRange(period instant.Period) (InstantRange, error) {
	value, err := temporalpostgres.NewInstantRange(period)
	return InstantRange{value: value}, err
}
func (r InstantRange) Period() (instant.Period, bool) { return r.value.Period() }
func (r InstantRange) Value() (driver.Value, error)   { return r.value.Value() }
func (r *InstantRange) Scan(source any) error {
	if r == nil {
		return temporal.ErrUnsupported
	}
	return r.value.Scan(source)
}

type DateRange struct{ value temporalpostgres.DateRange }

func NewDateRange(period dateperiod.Period) (DateRange, error) {
	value, err := temporalpostgres.NewDateRange(period)
	return DateRange{value: value}, err
}
func (r DateRange) Period() (dateperiod.Period, bool) { return r.value.Period() }
func (r DateRange) Value() (driver.Value, error)      { return r.value.Value() }
func (r *DateRange) Scan(source any) error {
	if r == nil {
		return temporal.ErrUnsupported
	}
	return r.value.Scan(source)
}
