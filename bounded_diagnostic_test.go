package temporal_test

import (
	"errors"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	calendar "github.com/faustbrian/go-calendar"
	temporal "github.com/faustbrian/go-temporal"
	newconfig "github.com/faustbrian/go-temporal/adapters/config"
	newpostgres "github.com/faustbrian/go-temporal/adapters/postgres"
	newwire "github.com/faustbrian/go-temporal/adapters/wire"
	"github.com/faustbrian/go-temporal/notation"
	oldconfig "github.com/faustbrian/go-temporal/temporalconfig"
	"github.com/faustbrian/go-temporal/temporalwire"
	"github.com/faustbrian/go-temporal/timeofday"
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
		"successor config": func() error {
			var value newconfig.InstantPeriod
			return value.UnmarshalText([]byte(secret))
		},
		"legacy config": func() error {
			var value oldconfig.InstantPeriod
			return value.UnmarshalText([]byte(secret))
		},
		"successor postgres": func() error {
			var value newpostgres.InstantRange
			return value.Scan(secret)
		},
		"successor wire": func() error {
			_, err := newwire.Unmarshal([]byte(`{"`+secret+`":true}`), temporal.Limits{})
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
	_, err := newwire.UnmarshalCollection(payload, temporal.Limits{InputPeriods: 1, ErrorBytes: 12})
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
			_, err := newwire.Unmarshal([]byte(
				`{"version":"temporal/v1","kind":"date-period","value":"[2026-02-30,2026-03-01)"}`,
			), temporal.Limits{})
			return err
		},
		"collection": func() error {
			_, err := newwire.UnmarshalCollection([]byte(
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

	_, err := (newwire.Document{}).Time(temporal.Limits{ErrorBytes: 1})
	if err == nil || len(err.Error()) > 1 || !errors.Is(err, temporal.ErrUnsupported) {
		t.Fatalf("Document.Time() error = %v", err)
	}

	document := newwire.CollectionDocument{
		Version: newwire.Version1,
		Kind:    newwire.KindInstantSet,
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
