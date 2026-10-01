package temporalconfig_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	temporal "github.com/faustbrian/go-temporal/v2"
	config "github.com/faustbrian/go-temporal/v2/adapters/config"
)

func TestConfigurationNilReceiversRejectWithoutPanic(t *testing.T) {
	for name, decode := range map[string]func([]byte) error{
		"instant":  (*config.InstantPeriod)(nil).UnmarshalText,
		"date":     (*config.DatePeriod)(nil).UnmarshalText,
		"daily":    (*config.DailyInterval)(nil).UnmarshalText,
		"duration": (*config.Duration)(nil).UnmarshalText,
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() != nil {
					t.Error("nil configuration receiver panicked")
				}
			}()
			if err := decode([]byte("ordinary-private-marker")); !errors.Is(err, temporal.ErrUnsupported) {
				t.Fatalf("nil receiver error = %v; want ErrUnsupported", err)
			}
		})
	}
}

func TestConfigurationTextAdmissionPreservesValueAndPrivateDiagnostics(t *testing.T) {
	// One fixed ordinary fixture is reused by all five wrappers; its inclusive
	// prefix is a slice, not a second boundary payload.
	text := []byte(strings.Repeat("a", 65537))
	for _, test := range []struct {
		name  string
		value interface {
			UnmarshalText([]byte) error
			MarshalText() ([]byte, error)
		}
		initial      string
		parseMessage string
	}{
		{"instant", &config.InstantPeriod{}, "[2026-01-02T03:04:05Z,2026-01-03T03:04:05Z)", "temporal: parse error: interval syntax"},
		{"date", &config.DatePeriod{}, "[2026-01-02,2026-01-03]", "temporal: parse error: interval syntax"},
		{"daily", &config.DailyInterval{}, "[08:00,09:00)", "temporal: parse error: interval syntax"},
		{"time", &config.Time{}, "08:00", "temporal: parse error: time syntax"},
		{"duration", &config.Duration{}, "PT1H", "temporal: parse error: duration syntax"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if err := test.value.UnmarshalText([]byte(test.initial)); err != nil {
				t.Fatal(err)
			}
			before, err := test.value.MarshalText()
			if err != nil || len(before) == 0 {
				t.Fatalf("initial value encoding = %q, %v", before, err)
			}
			err = test.value.UnmarshalText(text)
			var limit *temporal.LimitError
			if !errors.Is(err, temporal.ErrLimit) || !errors.As(err, &limit) ||
				limit.Field != "parse_bytes" || limit.Value != 65537 || limit.Max != 65536 ||
				err.Error() != "temporal: resource limit exceeded" {
				t.Fatalf("over-limit error = %v; want private fixed diagnostic with typed 65537/65536 byte limit", err)
			}
			after, err := test.value.MarshalText()
			if err != nil || !bytes.Equal(after, before) {
				t.Fatalf("over-limit rejection changed receiver: %q, %v", after, err)
			}
			err = test.value.UnmarshalText(text[:65536])
			if !errors.Is(err, temporal.ErrParse) || errors.Is(err, temporal.ErrLimit) || err.Error() != test.parseMessage {
				t.Fatalf("inclusive admission error = %v; want exact private syntax diagnostic, not ErrLimit", err)
			}
			after, err = test.value.MarshalText()
			if err != nil || !bytes.Equal(after, before) {
				t.Fatalf("inclusive malformed input changed receiver: %q, %v", after, err)
			}
		})
	}
}
