package temporalwire_test

import (
	"errors"
	"testing"

	calendar "github.com/faustbrian/go-calendar"
	temporal "github.com/faustbrian/go-temporal"
	temporalwire "github.com/faustbrian/go-temporal/adapters/wire"
	"github.com/faustbrian/go-temporal/timeofday"
)

func TestWireDocumentsAcceptEveryExactLimit(t *testing.T) {
	t.Parallel()

	encodedTime := timeofday.Noon().String()
	document, err := temporalwire.FromTime(timeofday.Noon(), temporal.Limits{FormatBytes: len(encodedTime)})
	if err != nil {
		t.Fatalf("FromTime(exact limit): %v", err)
	}
	payload, err := temporalwire.Marshal(document, temporal.Limits{})
	if err != nil {
		t.Fatalf("Marshal(): %v", err)
	}
	if _, err := temporalwire.Marshal(document, temporal.Limits{FormatBytes: len(payload)}); err != nil {
		t.Fatalf("Marshal(exact limit): %v", err)
	}
	if _, err := temporalwire.Unmarshal(payload, temporal.Limits{ParseBytes: len(payload)}); err != nil {
		t.Fatalf("Unmarshal(exact limit): %v", err)
	}
}

func TestDirectScalarDecodersBoundInvalidLimitsAndKinds(t *testing.T) {
	t.Parallel()

	invalid := temporal.Limits{ErrorBytes: -1}
	valid := temporal.Limits{ErrorBytes: 1}
	tests := map[string]func(temporal.Limits) error{
		"instant": func(limits temporal.Limits) error { _, err := (temporalwire.Document{}).Instant(limits); return err },
		"date":    func(limits temporal.Limits) error { _, err := (temporalwire.Document{}).Date(limits); return err },
		"daily": func(limits temporal.Limits) error {
			_, err := (temporalwire.Document{}).DailyInterval(limits)
			return err
		},
		"time":     func(limits temporal.Limits) error { _, err := (temporalwire.Document{}).Time(limits); return err },
		"duration": func(limits temporal.Limits) error { _, err := (temporalwire.Document{}).Duration(limits); return err },
	}
	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			var limitError *temporal.LimitError
			if err := run(invalid); !errors.As(err, &limitError) {
				t.Fatalf("invalid limits error = %v", err)
			}
			if err := run(valid); err == nil || len(err.Error()) > 1 || !errors.Is(err, temporal.ErrUnsupported) {
				t.Fatalf("kind error = %v", err)
			}
		})
	}
}

func TestDirectCollectionDecodersBoundEveryFailureStage(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		kind   temporalwire.Kind
		values []string
		run    func(temporalwire.CollectionDocument, temporal.Limits) error
	}{
		"instant": {temporalwire.KindInstantSet, []string{
			"[2026-01-01T00:00:00Z,2026-01-01T01:00:00Z)",
			"[2026-01-02T00:00:00Z,2026-01-02T01:00:00Z)",
		}, func(d temporalwire.CollectionDocument, limits temporal.Limits) error {
			_, err := d.InstantSet(limits)
			return err
		}},
		"date": {temporalwire.KindDateSet, []string{
			"[2026-01-01,2026-01-01]", "[2026-01-03,2026-01-03]",
		}, func(d temporalwire.CollectionDocument, limits temporal.Limits) error {
			_, err := d.DateSet(limits)
			return err
		}},
		"daily": {temporalwire.KindDailySet, []string{
			"[01:00,02:00)", "[03:00,04:00)",
		}, func(d temporalwire.CollectionDocument, limits temporal.Limits) error {
			_, err := d.DailySet(limits)
			return err
		}},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			document := temporalwire.CollectionDocument{Version: temporalwire.Version1, Kind: test.kind, Values: test.values}
			var limitError *temporal.LimitError
			if err := test.run(document, temporal.Limits{ErrorBytes: -1}); !errors.As(err, &limitError) {
				t.Fatalf("invalid limits error = %v", err)
			}
			if err := test.run(temporalwire.CollectionDocument{}, temporal.Limits{ErrorBytes: 1}); err == nil || len(err.Error()) > 1 || !errors.Is(err, temporal.ErrUnsupported) {
				t.Fatalf("kind error = %v", err)
			}
			if err := test.run(document, temporal.Limits{InputPeriods: 1, ErrorBytes: 1}); err == nil || len(err.Error()) > 1 || !errors.Is(err, temporal.ErrLimit) {
				t.Fatalf("input limit error = %v", err)
			}
			if err := test.run(document, temporal.Limits{InputPeriods: 2, OutputPeriods: 1, ErrorBytes: 1}); err == nil || len(err.Error()) > 1 || !errors.Is(err, temporal.ErrLimit) {
				t.Fatalf("output limit error = %v", err)
			}
		})
	}
}

func TestWireWrapperRetainsCalendarAndReversedClassifications(t *testing.T) {
	t.Parallel()

	_, err := temporalwire.Unmarshal([]byte(
		`{"version":"temporal/v1","kind":"date-period","value":"[2026-02-30,2026-03-01)"}`,
	), temporal.Limits{})
	if !errors.Is(err, calendar.ErrInvalidFormat) || !errors.Is(err, calendar.ErrInvalidDate) {
		t.Fatalf("calendar error = %v", err)
	}

	_, err = temporalwire.Unmarshal([]byte(
		`{"version":"temporal/v1","kind":"instant-period","value":"[2026-01-02T00:00:00Z,2026-01-01T00:00:00Z)"}`,
	), temporal.Limits{})
	if !errors.Is(err, temporal.ErrReversed) {
		t.Fatalf("reversed error = %v", err)
	}
}

func TestCollectionDiagnosticPreservesBoundedLimit(t *testing.T) {
	t.Parallel()

	payload := []byte(`{"version":"temporal/v1","kind":"instant-set","values":["a","b"]}`)
	_, err := temporalwire.UnmarshalCollection(payload, temporal.Limits{InputPeriods: 1, ErrorBytes: 12})
	if err == nil || len(err.Error()) > 12 || !errors.Is(err, temporal.ErrLimit) {
		t.Fatalf("UnmarshalCollection() error = %v", err)
	}
	var limitError *temporal.LimitError
	if !errors.As(err, &limitError) || limitError.Field != "input_periods" {
		t.Fatalf("LimitError = %#v", limitError)
	}
}

func TestWireCollectionsAcceptEveryExactLimit(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		document temporalwire.CollectionDocument
		decode   func(temporalwire.CollectionDocument, temporal.Limits) error
	}{
		"instant": {
			document: temporalwire.CollectionDocument{
				Version: temporalwire.Version1,
				Kind:    temporalwire.KindInstantSet,
				Values:  []string{"[2026-01-01T08:00:00Z,2026-01-01T17:00:00Z)"},
			},
			decode: func(document temporalwire.CollectionDocument, limits temporal.Limits) error {
				_, err := document.InstantSet(limits)
				return err
			},
		},
		"date": {
			document: temporalwire.CollectionDocument{
				Version: temporalwire.Version1,
				Kind:    temporalwire.KindDateSet,
				Values:  []string{"[2026-01-01,2026-01-02]"},
			},
			decode: func(document temporalwire.CollectionDocument, limits temporal.Limits) error {
				_, err := document.DateSet(limits)
				return err
			},
		},
		"daily": {
			document: temporalwire.CollectionDocument{
				Version: temporalwire.Version1,
				Kind:    temporalwire.KindDailySet,
				Values:  []string{"[08:00,17:00)"},
			},
			decode: func(document temporalwire.CollectionDocument, limits temporal.Limits) error {
				_, err := document.DailySet(limits)
				return err
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			limits := temporal.Limits{InputPeriods: 1}
			if err := test.decode(test.document, limits); err != nil {
				t.Fatalf("direct decode (exact input limit): %v", err)
			}

			payload, err := temporalwire.MarshalCollection(test.document, limits)
			if err != nil {
				t.Fatalf("MarshalCollection(exact input limit): %v", err)
			}
			if _, err := temporalwire.MarshalCollection(test.document, temporal.Limits{
				InputPeriods: 1,
				FormatBytes:  len(payload),
			}); err != nil {
				t.Fatalf("MarshalCollection(exact format limit): %v", err)
			}
			if _, err := temporalwire.UnmarshalCollection(payload, temporal.Limits{
				InputPeriods: 1,
				ParseBytes:   len(payload),
			}); err != nil {
				t.Fatalf("UnmarshalCollection(exact parse limit): %v", err)
			}
		})
	}
}
