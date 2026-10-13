// Package temporalwire provides retained wire adapters.
//
// Deprecated: use github.com/faustbrian/go-temporal/v2/adapters/wire. This
// package remains supported for the longer of 180 days after successor public
// availability and two subsequently published stable root-module minor
// releases.
package temporalwire

import (
	temporal "github.com/faustbrian/go-temporal/v2"
	adapter "github.com/faustbrian/go-temporal/v2/adapters/wire"
	"github.com/faustbrian/go-temporal/v2/dateperiod"
	"github.com/faustbrian/go-temporal/v2/instant"
	"github.com/faustbrian/go-temporal/v2/timeofday"
)

// Version1 names the retained temporal/v1 document schema.
const Version1 = adapter.Version1

// Kind selects a scalar or collection representation, checked during decoding.
type Kind string

const (
	// KindInstantPeriod selects an instant-period scalar.
	KindInstantPeriod Kind = "instant-period"
	// KindDatePeriod selects a civil-date-period scalar.
	KindDatePeriod Kind = "date-period"
	// KindDailyInterval selects a daily-interval scalar.
	KindDailyInterval Kind = "daily-interval"
	// KindTime selects a time-of-day scalar.
	KindTime Kind = "time-of-day"
	// KindDuration selects a fixed-duration scalar.
	KindDuration Kind = "fixed-duration"
	// KindInstantSet selects a normalized instant collection.
	KindInstantSet Kind = "instant-set"
	// KindDateSet selects a normalized civil-date collection.
	KindDateSet Kind = "date-set"
	// KindDailySet selects a normalized daily collection.
	KindDailySet Kind = "daily-set"
)

// Document is a caller-owned retained scalar envelope; constructing it does not
// validate its text. Its fields have the adapters/wire meanings.
type Document struct {
	// Version must equal Version1 for decoding.
	Version string `json:"version" yaml:"version" toml:"version"`
	// Kind must match the selected scalar accessor.
	Kind Kind `json:"kind" yaml:"kind" toml:"kind"`
	// Value holds interval or scalar notation, validated by the typed accessor.
	Value string `json:"value" yaml:"value" toml:"value"`
}

// FromInstant formats an instant value into a Version1 scalar envelope.
// It returns a zero document on failure; formatting is not revalidation against
// the selected parser profile. Marshal performs a separate admission check.
func FromInstant(value instant.Period, limits temporal.Limits) (Document, error) {
	result, err := adapter.FromInstant(value, limits)
	return fromAdapterDocument(result), err
}

// FromDate formats a civil-date value into a Version1 scalar envelope.
// It returns a zero document on failure; formatting is not revalidation against
// the selected parser profile. Marshal performs a separate admission check.
func FromDate(value dateperiod.Period, limits temporal.Limits) (Document, error) {
	result, err := adapter.FromDate(value, limits)
	return fromAdapterDocument(result), err
}

// FromDailyInterval formats a daily-interval value into a Version1 scalar envelope.
// It returns a zero document on failure; formatting is not revalidation against
// the selected parser profile. Marshal performs a separate admission check.
func FromDailyInterval(value timeofday.Interval, limits temporal.Limits) (Document, error) {
	result, err := adapter.FromDailyInterval(value, limits)
	return fromAdapterDocument(result), err
}

// FromTime formats a time-of-day value into a Version1 scalar envelope.
// It returns a zero document on failure; formatting is not revalidation against
// the selected parser profile. Marshal performs a separate admission check.
func FromTime(value timeofday.Time, limits temporal.Limits) (Document, error) {
	result, err := adapter.FromTime(value, limits)
	return fromAdapterDocument(result), err
}

// FromDuration formats a fixed-duration value into a Version1 scalar envelope.
// It returns a zero document on failure; formatting is not revalidation against
// the selected parser profile. Marshal performs a separate admission check.
func FromDuration(value timeofday.Duration, limits temporal.Limits) (Document, error) {
	result, err := adapter.FromDuration(value, limits)
	return fromAdapterDocument(result), err
}

// Instant checks the schema and kind, then parses Value using the supplied limits.
// It returns the core zero value on failure and does not mutate the document.
func (d Document) Instant(limits temporal.Limits) (instant.Period, error) {
	return d.adapter().Instant(limits)
}

// Date checks the schema and kind, then parses Value using the supplied limits.
// It returns the core zero value on failure and does not mutate the document.
func (d Document) Date(limits temporal.Limits) (dateperiod.Period, error) {
	return d.adapter().Date(limits)
}

// DailyInterval checks the schema and kind, then parses Value using the supplied limits.
// It returns the core zero value on failure and does not mutate the document.
func (d Document) DailyInterval(limits temporal.Limits) (timeofday.Interval, error) {
	return d.adapter().DailyInterval(limits)
}

// Time checks the schema and kind, then parses Value using the supplied limits.
// It returns the core zero value on failure and does not mutate the document.
func (d Document) Time(limits temporal.Limits) (timeofday.Time, error) {
	return d.adapter().Time(limits)
}

// Duration checks the schema and kind, then parses Value using the supplied limits.
// It returns the core zero value on failure and does not mutate the document.
func (d Document) Duration(limits temporal.Limits) (timeofday.Duration, error) {
	return d.adapter().Duration(limits)
}
func (d Document) adapter() adapter.Document {
	return adapter.Document{Version: d.Version, Kind: adapter.Kind(d.Kind), Value: d.Value}
}
func fromAdapterDocument(d adapter.Document) Document {
	return Document{Version: d.Version, Kind: Kind(d.Kind), Value: d.Value}
}

// Marshal validates a scalar and returns owned deterministic JSON.
// See adapters/wire for byte, depth, classification, and normalization rules.
func Marshal(document Document, limits temporal.Limits) ([]byte, error) {
	return adapter.Marshal(document.adapter(), limits)
}

// Unmarshal decodes one scalar JSON value with canonical adapter admission.
// See adapters/wire for byte, depth, classification, and normalization rules.
func Unmarshal(payload []byte, limits temporal.Limits) (Document, error) {
	result, err := adapter.Unmarshal(payload, limits)
	return fromAdapterDocument(result), err
}

// CollectionDocument is a caller-owned retained collection envelope. Values is
// borrowed during validation and encoding, so do not mutate it concurrently.
type CollectionDocument struct {
	// Version must equal Version1 for decoding.
	Version string `json:"version" yaml:"version" toml:"version"`
	// Kind must match the selected collection accessor.
	Kind Kind `json:"kind" yaml:"kind" toml:"kind"`
	// Values holds element notation; decoding normalizes order and overlap.
	Values []string `json:"values" yaml:"values" toml:"values"`
}

// FromInstantSet creates an owned outer Values slice for the instant set.
// The canonical adapter admits count, formats elements, and validates the envelope.
func FromInstantSet(set instant.Set, limits temporal.Limits) (CollectionDocument, error) {
	result, err := adapter.FromInstantSet(set, limits)
	return fromAdapterCollection(result), err
}

// FromDateSet creates an owned outer Values slice for the civil-date set.
// The canonical adapter admits count, formats elements, and validates the envelope.
func FromDateSet(set dateperiod.Set, limits temporal.Limits) (CollectionDocument, error) {
	result, err := adapter.FromDateSet(set, limits)
	return fromAdapterCollection(result), err
}

// FromDailySet creates an owned outer Values slice for the daily set.
// The canonical adapter admits count, formats elements, and validates the envelope.
func FromDailySet(set timeofday.IntervalSet, limits temporal.Limits) (CollectionDocument, error) {
	result, err := adapter.FromDailySet(set, limits)
	return fromAdapterCollection(result), err
}

// InstantSet borrows Values for the call, validates and normalizes it under limits,
// and returns a zero set on failure without mutating or retaining the slice.
func (d CollectionDocument) InstantSet(limits temporal.Limits) (instant.Set, error) {
	return d.adapter().InstantSet(limits)
}

// DateSet borrows Values for the call, validates and normalizes it under limits,
// and returns a zero set on failure without mutating or retaining the slice.
func (d CollectionDocument) DateSet(limits temporal.Limits) (dateperiod.Set, error) {
	return d.adapter().DateSet(limits)
}

// DailySet borrows Values for the call, validates and normalizes it under limits,
// and returns a zero set on failure without mutating or retaining the slice.
func (d CollectionDocument) DailySet(limits temporal.Limits) (timeofday.IntervalSet, error) {
	return d.adapter().DailySet(limits)
}
func (d CollectionDocument) adapter() adapter.CollectionDocument {
	// Canonical admission reads this borrowed view without retaining or mutating
	// Values. Copying here would allocate before its count and byte checks.
	return adapter.CollectionDocument{
		Version: d.Version, Kind: adapter.Kind(d.Kind), Values: d.Values,
	}
}
func fromAdapterCollection(d adapter.CollectionDocument) CollectionDocument {
	return CollectionDocument{
		Version: d.Version, Kind: Kind(d.Kind), Values: append([]string(nil), d.Values...),
	}
}

// MarshalCollection validates a collection and returns owned deterministic JSON.
// See adapters/wire for byte, depth, classification, and normalization rules.
func MarshalCollection(document CollectionDocument, limits temporal.Limits) ([]byte, error) {
	return adapter.MarshalCollection(document.adapter(), limits)
}

// UnmarshalCollection decodes and validates one collection JSON value.
// See adapters/wire for byte, depth, classification, and normalization rules.
func UnmarshalCollection(payload []byte, limits temporal.Limits) (CollectionDocument, error) {
	result, err := adapter.UnmarshalCollection(payload, limits)
	return fromAdapterCollection(result), err
}
