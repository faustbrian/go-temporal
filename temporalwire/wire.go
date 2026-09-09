// Package temporalwire provides retained wire adapters.
//
// Deprecated: use github.com/faustbrian/go-temporal/adapters/wire. This
// package remains supported for the longer of 180 days after successor public
// availability and two subsequently published stable root-module minor
// releases.
package temporalwire

import (
	temporal "github.com/faustbrian/go-temporal"
	adapter "github.com/faustbrian/go-temporal/adapters/wire"
	"github.com/faustbrian/go-temporal/dateperiod"
	"github.com/faustbrian/go-temporal/instant"
	"github.com/faustbrian/go-temporal/timeofday"
)

const Version1 = adapter.Version1

type Kind string

const (
	KindInstantPeriod Kind = "instant-period"
	KindDatePeriod    Kind = "date-period"
	KindDailyInterval Kind = "daily-interval"
	KindTime          Kind = "time-of-day"
	KindDuration      Kind = "fixed-duration"
	KindInstantSet    Kind = "instant-set"
	KindDateSet       Kind = "date-set"
	KindDailySet      Kind = "daily-set"
)

type Document struct {
	Version string `json:"version" yaml:"version" toml:"version"`
	Kind    Kind   `json:"kind" yaml:"kind" toml:"kind"`
	Value   string `json:"value" yaml:"value" toml:"value"`
}

func FromInstant(value instant.Period, limits temporal.Limits) (Document, error) {
	result, err := adapter.FromInstant(value, limits)
	return fromAdapterDocument(result), err
}
func FromDate(value dateperiod.Period, limits temporal.Limits) (Document, error) {
	result, err := adapter.FromDate(value, limits)
	return fromAdapterDocument(result), err
}
func FromDailyInterval(value timeofday.Interval, limits temporal.Limits) (Document, error) {
	result, err := adapter.FromDailyInterval(value, limits)
	return fromAdapterDocument(result), err
}
func FromTime(value timeofday.Time, limits temporal.Limits) (Document, error) {
	result, err := adapter.FromTime(value, limits)
	return fromAdapterDocument(result), err
}
func FromDuration(value timeofday.Duration, limits temporal.Limits) (Document, error) {
	result, err := adapter.FromDuration(value, limits)
	return fromAdapterDocument(result), err
}

func (d Document) Instant(limits temporal.Limits) (instant.Period, error) {
	return d.adapter().Instant(limits)
}
func (d Document) Date(limits temporal.Limits) (dateperiod.Period, error) {
	return d.adapter().Date(limits)
}
func (d Document) DailyInterval(limits temporal.Limits) (timeofday.Interval, error) {
	return d.adapter().DailyInterval(limits)
}
func (d Document) Time(limits temporal.Limits) (timeofday.Time, error) {
	return d.adapter().Time(limits)
}
func (d Document) Duration(limits temporal.Limits) (timeofday.Duration, error) {
	return d.adapter().Duration(limits)
}
func (d Document) adapter() adapter.Document {
	return adapter.Document{Version: d.Version, Kind: adapter.Kind(d.Kind), Value: d.Value}
}
func fromAdapterDocument(d adapter.Document) Document {
	return Document{Version: d.Version, Kind: Kind(d.Kind), Value: d.Value}
}

func Marshal(document Document, limits temporal.Limits) ([]byte, error) {
	return adapter.Marshal(document.adapter(), limits)
}
func Unmarshal(payload []byte, limits temporal.Limits) (Document, error) {
	result, err := adapter.Unmarshal(payload, limits)
	return fromAdapterDocument(result), err
}

type CollectionDocument struct {
	Version string   `json:"version" yaml:"version" toml:"version"`
	Kind    Kind     `json:"kind" yaml:"kind" toml:"kind"`
	Values  []string `json:"values" yaml:"values" toml:"values"`
}

func FromInstantSet(set instant.Set, limits temporal.Limits) (CollectionDocument, error) {
	result, err := adapter.FromInstantSet(set, limits)
	return fromAdapterCollection(result), err
}
func FromDateSet(set dateperiod.Set, limits temporal.Limits) (CollectionDocument, error) {
	result, err := adapter.FromDateSet(set, limits)
	return fromAdapterCollection(result), err
}
func FromDailySet(set timeofday.IntervalSet, limits temporal.Limits) (CollectionDocument, error) {
	result, err := adapter.FromDailySet(set, limits)
	return fromAdapterCollection(result), err
}
func (d CollectionDocument) InstantSet(limits temporal.Limits) (instant.Set, error) {
	return d.adapter().InstantSet(limits)
}
func (d CollectionDocument) DateSet(limits temporal.Limits) (dateperiod.Set, error) {
	return d.adapter().DateSet(limits)
}
func (d CollectionDocument) DailySet(limits temporal.Limits) (timeofday.IntervalSet, error) {
	return d.adapter().DailySet(limits)
}
func (d CollectionDocument) adapter() adapter.CollectionDocument {
	return adapter.CollectionDocument{
		Version: d.Version, Kind: adapter.Kind(d.Kind), Values: append([]string(nil), d.Values...),
	}
}
func fromAdapterCollection(d adapter.CollectionDocument) CollectionDocument {
	return CollectionDocument{
		Version: d.Version, Kind: Kind(d.Kind), Values: append([]string(nil), d.Values...),
	}
}
func MarshalCollection(document CollectionDocument, limits temporal.Limits) ([]byte, error) {
	return adapter.MarshalCollection(document.adapter(), limits)
}
func UnmarshalCollection(payload []byte, limits temporal.Limits) (CollectionDocument, error) {
	result, err := adapter.UnmarshalCollection(payload, limits)
	return fromAdapterCollection(result), err
}
