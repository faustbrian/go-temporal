package temporalwire

import (
	"bytes"
	"encoding/json"
	"io"
	"unicode/utf8"

	temporal "github.com/faustbrian/go-temporal"
	"github.com/faustbrian/go-temporal/dateperiod"
	"github.com/faustbrian/go-temporal/instant"
	"github.com/faustbrian/go-temporal/internal/diagnostic"
	"github.com/faustbrian/go-temporal/notation"
	"github.com/faustbrian/go-temporal/timeofday"
)

// CollectionDocument is a stable wire envelope for one normalized interval
// set. Values are canonical ISO 80000 elements in normalized order.
type CollectionDocument struct {
	Version string   `json:"version" yaml:"version" toml:"version"`
	Kind    Kind     `json:"kind" yaml:"kind" toml:"kind"`
	Values  []string `json:"values" yaml:"values" toml:"values"`
}

// FromInstantSet constructs a versioned normalized instant-set document.
func FromInstantSet(set instant.Set, limits temporal.Limits) (CollectionDocument, error) {
	values := make([]string, 0, set.Len())
	for _, period := range set.Periods() {
		encoded, err := notation.FormatInstant(period, notation.ISO80000, limits)
		if err != nil {
			return CollectionDocument{}, err
		}
		values = append(values, encoded)
	}
	return newCollection(KindInstantSet, values, limits)
}

// FromDateSet constructs a versioned normalized civil-date-set document.
func FromDateSet(set dateperiod.Set, limits temporal.Limits) (CollectionDocument, error) {
	values := make([]string, 0, set.Len())
	for _, period := range set.Periods() {
		encoded, err := notation.FormatDate(period, notation.ISO80000, limits)
		if err != nil {
			return CollectionDocument{}, err
		}
		values = append(values, encoded)
	}
	return newCollection(KindDateSet, values, limits)
}

// FromDailySet constructs a versioned normalized daily-set document.
func FromDailySet(set timeofday.IntervalSet, limits temporal.Limits) (CollectionDocument, error) {
	values := make([]string, 0, set.Len())
	for _, interval := range set.Intervals() {
		encoded, err := notation.FormatDailyInterval(interval, notation.ISO80000, limits)
		if err != nil {
			return CollectionDocument{}, err
		}
		values = append(values, encoded)
	}
	return newCollection(KindDailySet, values, limits)
}

func newCollection(kind Kind, values []string, limits temporal.Limits) (CollectionDocument, error) {
	document := CollectionDocument{Version: Version1, Kind: kind, Values: append([]string(nil), values...)}
	if err := document.validate(limits); err != nil {
		return CollectionDocument{}, err
	}
	return document, nil
}

// InstantSet decodes and normalizes an instant-set document.
func (d CollectionDocument) InstantSet(limits temporal.Limits) (instant.Set, error) {
	limits = limits.Resolve()
	if err := limits.Validate(); err != nil {
		return instant.Set{}, err
	}
	if err := d.expect(KindInstantSet); err != nil {
		return instant.Set{}, boundedWireError(limits, "collection document kind", err)
	}
	if len(d.Values) > limits.InputPeriods {
		cause := &temporal.LimitError{Field: "input_periods", Value: len(d.Values), Max: limits.InputPeriods}
		return instant.Set{}, boundedWireError(limits, "collection document size", cause)
	}
	periods := make([]instant.Period, 0, len(d.Values))
	for _, value := range d.Values {
		period, err := notation.ParseInstant(value, notation.ISO80000, limits)
		if err != nil {
			return instant.Set{}, err
		}
		periods = append(periods, period)
	}
	set, err := instant.NewSet(limits, periods...)
	if err != nil {
		return instant.Set{}, boundedWireError(limits, "collection document value", err)
	}
	return set, nil
}

// DateSet decodes and normalizes a civil-date-set document.
func (d CollectionDocument) DateSet(limits temporal.Limits) (dateperiod.Set, error) {
	limits = limits.Resolve()
	if err := limits.Validate(); err != nil {
		return dateperiod.Set{}, err
	}
	if err := d.expect(KindDateSet); err != nil {
		return dateperiod.Set{}, boundedWireError(limits, "collection document kind", err)
	}
	if len(d.Values) > limits.InputPeriods {
		cause := &temporal.LimitError{Field: "input_periods", Value: len(d.Values), Max: limits.InputPeriods}
		return dateperiod.Set{}, boundedWireError(limits, "collection document size", cause)
	}
	periods := make([]dateperiod.Period, 0, len(d.Values))
	for _, value := range d.Values {
		period, err := notation.ParseDate(value, notation.ISO80000, limits)
		if err != nil {
			return dateperiod.Set{}, err
		}
		periods = append(periods, period)
	}
	set, err := dateperiod.NewSet(limits, periods...)
	if err != nil {
		return dateperiod.Set{}, boundedWireError(limits, "collection document value", err)
	}
	return set, nil
}

// DailySet decodes and normalizes a daily interval-set document.
func (d CollectionDocument) DailySet(limits temporal.Limits) (timeofday.IntervalSet, error) {
	limits = limits.Resolve()
	if err := limits.Validate(); err != nil {
		return timeofday.IntervalSet{}, err
	}
	if err := d.expect(KindDailySet); err != nil {
		return timeofday.IntervalSet{}, boundedWireError(limits, "collection document kind", err)
	}
	if len(d.Values) > limits.InputPeriods {
		cause := &temporal.LimitError{Field: "input_periods", Value: len(d.Values), Max: limits.InputPeriods}
		return timeofday.IntervalSet{}, boundedWireError(limits, "collection document size", cause)
	}
	intervals := make([]timeofday.Interval, 0, len(d.Values))
	for _, value := range d.Values {
		interval, err := notation.ParseDailyInterval(value, notation.ISO80000, limits)
		if err != nil {
			return timeofday.IntervalSet{}, err
		}
		intervals = append(intervals, interval)
	}
	set, err := timeofday.NewIntervalSet(limits, intervals...)
	if err != nil {
		return timeofday.IntervalSet{}, boundedWireError(limits, "collection document value", err)
	}
	return set, nil
}

func (d CollectionDocument) expect(kind Kind) error {
	if d.Version != Version1 || d.Kind != kind {
		return temporal.ErrUnsupported
	}
	return nil
}

func (d CollectionDocument) validate(limits temporal.Limits) error {
	limits = limits.Resolve()
	if err := limits.Validate(); err != nil {
		return err
	}
	if len(d.Values) > limits.InputPeriods {
		return &temporal.LimitError{Field: "input_periods", Value: len(d.Values), Max: limits.InputPeriods}
	}
	switch d.Kind {
	case KindInstantSet:
		_, err := d.InstantSet(limits)
		return err
	case KindDateSet:
		_, err := d.DateSet(limits)
		return err
	case KindDailySet:
		_, err := d.DailySet(limits)
		return err
	default:
		return temporal.ErrUnsupported
	}
}

// MarshalCollection returns deterministic JSON for a valid collection.
func MarshalCollection(document CollectionDocument, limits temporal.Limits) ([]byte, error) {
	limits = limits.Resolve()
	if err := document.validate(limits); err != nil {
		return nil, err
	}
	document.Values = append([]string(nil), document.Values...)
	payload, _ := json.Marshal(document)
	if len(payload) > limits.FormatBytes {
		return nil, &temporal.LimitError{Field: "format_bytes", Value: len(payload), Max: limits.FormatBytes}
	}
	return payload, nil
}

// UnmarshalCollection strictly decodes exactly one collection document.
func UnmarshalCollection(payload []byte, limits temporal.Limits) (CollectionDocument, error) {
	limits = limits.Resolve()
	if err := limits.Validate(); err != nil {
		return CollectionDocument{}, err
	}
	if len(payload) > limits.ParseBytes {
		cause := &temporal.LimitError{Field: "parse_bytes", Value: len(payload), Max: limits.ParseBytes}
		return CollectionDocument{}, diagnostic.New(limits.ErrorBytes, temporal.ErrLimit.Error(), temporal.ErrLimit, cause)
	}
	if !utf8.Valid(payload) {
		return CollectionDocument{}, diagnostic.New(limits.ErrorBytes, "temporal: parse error: collection document syntax", temporal.ErrParse)
	}
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	var document CollectionDocument
	if err := decoder.Decode(&document); err != nil {
		return CollectionDocument{}, diagnostic.New(limits.ErrorBytes, "temporal: parse error: collection document syntax", temporal.ErrParse)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return CollectionDocument{}, diagnostic.New(limits.ErrorBytes, "temporal: parse error: trailing collection document", temporal.ErrParse)
	}
	if err := document.validate(limits); err != nil {
		return CollectionDocument{}, boundedWireError(limits, "collection document value", err)
	}
	document.Values = append([]string(nil), document.Values...)
	return document, nil
}
