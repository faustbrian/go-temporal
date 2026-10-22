package temporal_test

import (
	"errors"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	calendar "github.com/faustbrian/go-calendar/v2"
	temporal "github.com/faustbrian/go-temporal/v3"
	temporalconfig "github.com/faustbrian/go-temporal/v3/adapters/config"
	temporalpostgres "github.com/faustbrian/go-temporal/v3/adapters/postgres"
	temporalwire "github.com/faustbrian/go-temporal/v3/adapters/wire"
	"github.com/faustbrian/go-temporal/v3/notation"
	"github.com/faustbrian/go-temporal/v3/timeofday"
)

func TestHostileDiagnosticsAreBoundedAndSanitized(t *testing.T) {
	t.Parallel()

	secret := strings.Repeat("attacker-controlled-", 80)
	tests := map[string]struct {
		limit int
		run   func() error
		is    error
	}{
		"bounds": {
			limit: temporal.DefaultLimits().ErrorBytes,
			run: func() error {
				var value temporal.Bounds
				return value.UnmarshalText([]byte(secret))
			},
			is: temporal.ErrBounds,
		},
		"instant": {
			limit: 20,
			run: func() error {
				_, err := notation.ParseInstant(
					secret+"/2026-01-01T00:00:00Z", notation.ISO8601,
					temporal.Limits{ErrorBytes: 20},
				)
				return err
			},
			is: temporal.ErrParse,
		},
		"wire scalar member": {
			limit: 24,
			run: func() error {
				_, err := temporalwire.Unmarshal(
					[]byte(`{"version":"temporal/v1","`+secret+`":true}`),
					temporal.Limits{ErrorBytes: 24},
				)
				return err
			},
			is: temporal.ErrParse,
		},
		"wire collection member": {
			limit: 24,
			run: func() error {
				_, err := temporalwire.UnmarshalCollection(
					[]byte(`{"version":"temporal/v1","`+secret+`":true}`),
					temporal.Limits{ErrorBytes: 24},
				)
				return err
			},
			is: temporal.ErrParse,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.run()
			if err == nil {
				t.Fatal("expected error")
			}
			if len(err.Error()) > test.limit {
				t.Fatalf("error length = %d, want <= %d: %q", len(err.Error()), test.limit, err)
			}
			if !utf8.ValidString(err.Error()) {
				t.Fatalf("error is not valid UTF-8: %q", err)
			}
			if strings.Contains(err.Error(), "attacker-controlled") {
				t.Fatalf("error retained hostile input: %q", err)
			}
			if !errors.Is(err, test.is) {
				t.Fatalf("errors.Is(%v) = false for %v", err, test.is)
			}
			var parseError *time.ParseError
			if errors.As(err, &parseError) {
				t.Fatalf("error retained unsafe parser cause: %#v", parseError)
			}
		})
	}
}

func TestFractionalTimeInvalidValueUsesBoundedDiagnostic(t *testing.T) {
	t.Parallel()

	_, err := timeofday.Parse("99:00:00.1", temporal.Limits{ErrorBytes: 5})
	if err == nil || len(err.Error()) > 5 || !errors.Is(err, temporal.ErrInvalidTime) {
		t.Fatalf("err = %v", err)
	}
}

func TestAdapterHostileTextBoundariesUseTheSameDiagnosticContract(t *testing.T) {
	t.Parallel()

	secret := strings.Repeat("private-value-", 100)
	checks := map[string]func() error{
		"config adapter": func() error {
			var value temporalconfig.InstantPeriod
			return value.UnmarshalText([]byte(secret))
		},
		"postgres adapter": func() error {
			var value temporalpostgres.InstantRange
			return value.Scan(secret)
		},
		"wire adapter": func() error {
			_, err := temporalwire.Unmarshal([]byte(`{"`+secret+`":true}`), temporal.Limits{})
			return err
		},
	}
	for name, check := range checks {
		t.Run(name, func(t *testing.T) {
			err := check()
			if err == nil || len(err.Error()) > temporal.DefaultLimits().ErrorBytes ||
				strings.Contains(err.Error(), "private-value") {
				t.Fatalf("err = %v", err)
			}
		})
	}
}

func TestTinyDiagnosticBudgetsPreserveClassification(t *testing.T) {
	t.Parallel()

	for budget := 1; budget <= len("temporal: parse error: timestamp syntax"); budget++ {
		_, err := notation.ParseInstant("hostile/2026-01-01T00:00:00Z", notation.ISO8601,
			temporal.Limits{ErrorBytes: budget})
		if err == nil || !errors.Is(err, temporal.ErrParse) {
			t.Fatalf("budget %d: err = %v", budget, err)
		}
		if len(err.Error()) > budget || !utf8.ValidString(err.Error()) {
			t.Fatalf("budget %d: invalid bounded text %q", budget, err)
		}
	}
}

func TestParseLimitKeepsSafeStructuredCause(t *testing.T) {
	t.Parallel()

	_, err := notation.ParseInstant(strings.Repeat("x", 11), notation.ISO8601,
		temporal.Limits{ParseBytes: 10, ErrorBytes: 8})
	if err == nil || len(err.Error()) > 8 || !errors.Is(err, temporal.ErrLimit) {
		t.Fatalf("err = %v", err)
	}
	var limitError *temporal.LimitError
	if !errors.As(err, &limitError) {
		t.Fatalf("errors.As(*LimitError) = false: %v", err)
	}
	if limitError.Field != "parse_bytes" || limitError.Value != 11 || limitError.Max != 10 {
		t.Fatalf("LimitError = %#v", limitError)
	}
}

func TestWireCollectionLimitKeepsSafeStructuredCause(t *testing.T) {
	t.Parallel()

	payload := []byte(`{"version":"temporal/v1","kind":"instant-set","values":["a","b"]}`)
	_, err := temporalwire.UnmarshalCollection(payload, temporal.Limits{InputPeriods: 1, ErrorBytes: 12})
	if err == nil || len(err.Error()) > 12 || !errors.Is(err, temporal.ErrLimit) {
		t.Fatalf("err = %v", err)
	}
	var limitError *temporal.LimitError
	if !errors.As(err, &limitError) || limitError.Field != "input_periods" {
		t.Fatalf("LimitError = %#v", limitError)
	}
}

func TestDateDiagnosticKeepsSafeCalendarClassification(t *testing.T) {
	t.Parallel()

	_, err := notation.ParseDate("2026-02-30/2026-03-01", notation.ISO8601,
		temporal.Limits{ErrorBytes: 16})
	if err == nil || len(err.Error()) > 16 || !errors.Is(err, temporal.ErrParse) ||
		!errors.Is(err, calendar.ErrInvalidFormat) || !errors.Is(err, calendar.ErrInvalidDate) {
		t.Fatalf("err = %v", err)
	}
}

func TestReversedDiagnosticsPreserveClassification(t *testing.T) {
	t.Parallel()

	for name, run := range map[string]func() error{
		"instant": func() error {
			_, err := notation.ParseInstant(
				"2026-01-02T00:00:00Z/2026-01-01T00:00:00Z", notation.ISO8601,
				temporal.Limits{ErrorBytes: 12},
			)
			return err
		},
		"date": func() error {
			_, err := notation.ParseDate("2026-01-02/2026-01-01", notation.ISO8601,
				temporal.Limits{ErrorBytes: 12})
			return err
		},
	} {
		t.Run(name, func(t *testing.T) {
			err := run()
			if err == nil || len(err.Error()) > 12 || !errors.Is(err, temporal.ErrReversed) {
				t.Fatalf("err = %v", err)
			}
		})
	}
}

func TestWireDateDiagnosticsPreserveCalendarClassification(t *testing.T) {
	t.Parallel()

	for name, run := range map[string]func() error{
		"scalar": func() error {
			_, err := temporalwire.Unmarshal([]byte(
				`{"version":"temporal/v1","kind":"date-period","value":"[2026-02-30,2026-03-01)"}`,
			), temporal.Limits{})
			return err
		},
		"collection": func() error {
			_, err := temporalwire.UnmarshalCollection([]byte(
				`{"version":"temporal/v1","kind":"date-set","values":["[2026-02-30,2026-03-01)"]}`,
			), temporal.Limits{})
			return err
		},
	} {
		t.Run(name, func(t *testing.T) {
			err := run()
			if err == nil || !errors.Is(err, calendar.ErrInvalidFormat) ||
				!errors.Is(err, calendar.ErrInvalidDate) {
				t.Fatalf("err = %v", err)
			}
		})
	}
}

func TestDirectWireDecodersBoundEveryTerminalError(t *testing.T) {
	t.Parallel()

	_, err := (temporalwire.Document{}).Time(temporal.Limits{ErrorBytes: 1})
	if err == nil || len(err.Error()) > 1 || !errors.Is(err, temporal.ErrUnsupported) {
		t.Fatalf("Document.Time() error = %v", err)
	}

	document := temporalwire.CollectionDocument{
		Version: temporalwire.Version1,
		Kind:    temporalwire.KindInstantSet,
		Values: []string{
			"[2026-01-01T00:00:00Z,2026-01-01T01:00:00Z)",
			"[2026-01-02T00:00:00Z,2026-01-02T01:00:00Z)",
		},
	}
	_, err = document.InstantSet(temporal.Limits{InputPeriods: 1, ErrorBytes: 1})
	if err == nil || len(err.Error()) > 1 || !errors.Is(err, temporal.ErrLimit) {
		t.Fatalf("CollectionDocument.InstantSet() error = %v", err)
	}
	var limitError *temporal.LimitError
	if !errors.As(err, &limitError) || limitError.Field != "input_periods" {
		t.Fatalf("LimitError = %#v", limitError)
	}
}
