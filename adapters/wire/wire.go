// Package temporalwire provides versioned, format-neutral documents for
// encoding temporal values through wire or the standard JSON package.
package temporalwire

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"unicode/utf8"

	calendar "github.com/faustbrian/go-calendar/v2"
	temporal "github.com/faustbrian/go-temporal/v2"
	"github.com/faustbrian/go-temporal/v2/dateperiod"
	"github.com/faustbrian/go-temporal/v2/instant"
	"github.com/faustbrian/go-temporal/v2/internal/diagnostic"
	"github.com/faustbrian/go-temporal/v2/notation"
	"github.com/faustbrian/go-temporal/v2/timeofday"
)

// Version1 is the stable initial document schema identifier.
const Version1 = "temporal/v1"

// Kind identifies the temporal value encoded by a document.
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

// Document is a format-neutral stable wire representation. Value is canonical
// ISO 80000 notation for intervals and strict ISO text for scalar values.
type Document struct {
	Version string `json:"version" yaml:"version" toml:"version"`
	Kind    Kind   `json:"kind" yaml:"kind" toml:"kind"`
	Value   string `json:"value" yaml:"value" toml:"value"`
}

func FromInstant(value instant.Period, limits temporal.Limits) (Document, error) {
	encoded, err := notation.FormatInstant(value, notation.ISO80000, limits)
	return fromEncoded(KindInstantPeriod, encoded, err)
}

func FromDate(value dateperiod.Period, limits temporal.Limits) (Document, error) {
	encoded, err := notation.FormatDate(value, notation.ISO80000, limits)
	return fromEncoded(KindDatePeriod, encoded, err)
}

func FromDailyInterval(value timeofday.Interval, limits temporal.Limits) (Document, error) {
	encoded, err := notation.FormatDailyInterval(value, notation.ISO80000, limits)
	return fromEncoded(KindDailyInterval, encoded, err)
}

func FromTime(value timeofday.Time, limits temporal.Limits) (Document, error) {
	limits = limits.Resolve()
	if err := limits.Validate(); err != nil {
		return Document{}, err
	}
	encoded := value.String()
	if len(encoded) > limits.FormatBytes {
		return Document{}, &temporal.LimitError{Field: "format_bytes", Value: len(encoded), Max: limits.FormatBytes}
	}
	return Document{Version: Version1, Kind: KindTime, Value: encoded}, nil
}

func FromDuration(value timeofday.Duration, limits temporal.Limits) (Document, error) {
	encoded, err := notation.FormatDuration(value, limits)
	return fromEncoded(KindDuration, encoded, err)
}

func fromEncoded(kind Kind, encoded string, err error) (Document, error) {
	if err != nil {
		return Document{}, err
	}
	return Document{Version: Version1, Kind: kind, Value: encoded}, nil
}

func (d Document) Instant(limits temporal.Limits) (instant.Period, error) {
	limits = limits.Resolve()
	if err := limits.Validate(); err != nil {
		return instant.Period{}, err
	}
	if err := d.expect(KindInstantPeriod); err != nil {
		return instant.Period{}, boundedWireError(limits, "scalar document kind", err)
	}
	return notation.ParseInstant(d.Value, notation.ISO80000, limits)
}

func (d Document) Date(limits temporal.Limits) (dateperiod.Period, error) {
	limits = limits.Resolve()
	if err := limits.Validate(); err != nil {
		return dateperiod.Period{}, err
	}
	if err := d.expect(KindDatePeriod); err != nil {
		return dateperiod.Period{}, boundedWireError(limits, "scalar document kind", err)
	}
	return notation.ParseDate(d.Value, notation.ISO80000, limits)
}

func (d Document) DailyInterval(limits temporal.Limits) (timeofday.Interval, error) {
	limits = limits.Resolve()
	if err := limits.Validate(); err != nil {
		return timeofday.Interval{}, err
	}
	if err := d.expect(KindDailyInterval); err != nil {
		return timeofday.Interval{}, boundedWireError(limits, "scalar document kind", err)
	}
	return notation.ParseDailyInterval(d.Value, notation.ISO80000, limits)
}

func (d Document) Time(limits temporal.Limits) (timeofday.Time, error) {
	limits = limits.Resolve()
	if err := limits.Validate(); err != nil {
		return timeofday.Time{}, err
	}
	if err := d.expect(KindTime); err != nil {
		return timeofday.Time{}, boundedWireError(limits, "scalar document kind", err)
	}
	return timeofday.Parse(d.Value, limits)
}

func (d Document) Duration(limits temporal.Limits) (timeofday.Duration, error) {
	limits = limits.Resolve()
	if err := limits.Validate(); err != nil {
		return timeofday.Duration{}, err
	}
	if err := d.expect(KindDuration); err != nil {
		return timeofday.Duration{}, boundedWireError(limits, "scalar document kind", err)
	}
	return notation.ParseDuration(d.Value, limits)
}

func (d Document) expect(kind Kind) error {
	if d.Version != Version1 || d.Kind != kind {
		return temporal.ErrUnsupported
	}
	return nil
}

func (d Document) validate(limits temporal.Limits) error {
	switch d.Kind {
	case KindInstantPeriod:
		_, err := d.Instant(limits)
		return err
	case KindDatePeriod:
		_, err := d.Date(limits)
		return err
	case KindDailyInterval:
		_, err := d.DailyInterval(limits)
		return err
	case KindTime:
		_, err := d.Time(limits)
		return err
	case KindDuration:
		_, err := d.Duration(limits)
		return err
	default:
		return temporal.ErrUnsupported
	}
}

// Marshal returns deterministic JSON for a valid versioned document.
func Marshal(document Document, limits temporal.Limits) ([]byte, error) {
	limits = limits.Resolve()
	if err := limits.Validate(); err != nil {
		return nil, err
	}
	if err := document.validate(limits); err != nil {
		return nil, err
	}
	if size, exceeds := scalarDocumentSize(document, limits.FormatBytes); exceeds {
		return nil, &temporal.LimitError{Field: "format_bytes", Value: size, Max: limits.FormatBytes}
	}
	payload, _ := json.Marshal(document)
	if len(payload) > limits.FormatBytes {
		return nil, &temporal.LimitError{Field: "format_bytes", Value: len(payload), Max: limits.FormatBytes}
	}
	return payload, nil
}

func scalarDocumentSize(document Document, maximum int) (int, bool) {
	size := len(`{"version":`) + jsonStringSize(document.Version) +
		len(`,"kind":`) + jsonStringSize(string(document.Kind)) +
		len(`,"value":`) + jsonStringSize(document.Value) + 1
	return boundedJSONSize(size, maximum)
}

func jsonStringSize(value string) int {
	// Count admitted UTF-8 strings without allocating their JSON representation.
	// Document validation has already rejected invalid UTF-8. The replacement
	// case remains a conservative bound for an invalid byte.
	size := 2
	for len(value) > 0 {
		r, width := utf8.DecodeRuneInString(value)
		value = value[width:]
		switch {
		case r == '"' || r == '\\' || r == '\b' || r == '\f' || r == '\n' || r == '\r' || r == '\t':
			size += 2
		case r < 0x20 || r == '<' || r == '>' || r == '&' || r == '\u2028' || r == '\u2029' || (r == utf8.RuneError && width == 1):
			size += 6
		default:
			size += width
		}
	}
	return size
}

func boundedJSONSize(size, maximum int) (int, bool) {
	if size > maximum {
		return size, true
	}
	return size, false
}

// Unmarshal strictly decodes exactly one versioned JSON document.
func Unmarshal(payload []byte, limits temporal.Limits) (Document, error) {
	limits = limits.Resolve()
	if err := limits.Validate(); err != nil {
		return Document{}, err
	}
	if len(payload) > limits.ParseBytes {
		cause := &temporal.LimitError{Field: "parse_bytes", Value: len(payload), Max: limits.ParseBytes}
		return Document{}, diagnostic.New(limits.ErrorBytes, temporal.ErrLimit.Error(), temporal.ErrLimit, cause)
	}
	if !utf8.Valid(payload) {
		return Document{}, diagnostic.New(limits.ErrorBytes, "temporal: parse error: scalar document syntax", temporal.ErrParse)
	}
	if err := validateJSONStructure(payload, limits.ParserDepth); err != nil {
		return Document{}, boundedWireError(limits, "scalar document syntax", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	var document Document
	if err := decoder.Decode(&document); err != nil {
		return Document{}, diagnostic.New(limits.ErrorBytes, "temporal: parse error: scalar document syntax", temporal.ErrParse)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return Document{}, diagnostic.New(limits.ErrorBytes, "temporal: parse error: trailing document", temporal.ErrParse)
	}
	if err := document.validate(limits); err != nil {
		return Document{}, boundedWireError(limits, "scalar document value", err)
	}
	return document, nil
}

func validateJSONStructure(payload []byte, maxDepth int) error {
	decoder := json.NewDecoder(bytes.NewReader(payload))
	token, err := decoder.Token()
	if err != nil {
		return temporal.ErrParse
	}
	if err := validateJSONValue(decoder, token, 1, maxDepth); err != nil {
		return err
	}
	if _, err := decoder.Token(); err != io.EOF {
		return temporal.ErrParse
	}
	return nil
}

func validateJSONValue(decoder *json.Decoder, token json.Token, depth, maxDepth int) error {
	delimiter, compound := token.(json.Delim)
	if !compound {
		return nil
	}
	if depth > maxDepth {
		return &temporal.LimitError{Field: "parser_depth", Value: depth, Max: maxDepth}
	}

	switch delimiter {
	case '{':
		keys := make(map[string]struct{})
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return temporal.ErrParse
			}
			key, ok := keyToken.(string)
			if !ok {
				return temporal.ErrParse
			}
			if _, exists := keys[key]; exists {
				return temporal.ErrParse
			}
			keys[key] = struct{}{}
			valueToken, err := decoder.Token()
			if err != nil {
				return temporal.ErrParse
			}
			if err := validateJSONValue(decoder, valueToken, depth+1, maxDepth); err != nil {
				return err
			}
		}
	case '[':
		for decoder.More() {
			valueToken, err := decoder.Token()
			if err != nil {
				return temporal.ErrParse
			}
			if err := validateJSONValue(decoder, valueToken, depth+1, maxDepth); err != nil {
				return err
			}
		}
	default:
		return temporal.ErrParse
	}

	closing, err := decoder.Token()
	if err != nil {
		return temporal.ErrParse
	}
	expected := json.Delim('}')
	if delimiter == '[' {
		expected = ']'
	}
	if closing != expected {
		return temporal.ErrParse
	}
	return nil
}

func boundedWireError(limits temporal.Limits, stage string, cause error) error {
	causes := make([]error, 0, 4)
	for _, sentinel := range []error{temporal.ErrBounds, temporal.ErrPrecision, temporal.ErrInvalidTime,
		temporal.ErrUnsupported, temporal.ErrOverflow, temporal.ErrLimit, temporal.ErrParse,
		temporal.ErrReversed} {
		if errors.Is(cause, sentinel) {
			causes = append(causes, sentinel)
		}
	}
	var limitError *temporal.LimitError
	if errors.As(cause, &limitError) {
		causes = append(causes, limitError)
	}
	for _, sentinel := range []error{calendar.ErrInvalidFormat, calendar.ErrInvalidDate} {
		if errors.Is(cause, sentinel) {
			causes = append(causes, sentinel)
		}
	}
	return diagnostic.New(limits.ErrorBytes, "temporal: parse error: "+stage, causes...)
}
