package notation_test

import (
	"errors"
	"regexp"
	"testing"
	"time"
	"unicode/utf8"

	temporal "github.com/faustbrian/go-temporal/v3"
	"github.com/faustbrian/go-temporal/v3/notation"
)

func TestInstantEndpointSelectedProfileRejectsFallbackExtensions(t *testing.T) {
	validEnd := "2026-10-15T12:00:00Z"
	for _, start := range []string{
		"2026-10-13T1:00:00Z",
		"2026-10-13T1:00:00.123Z",
		"2026-10-13T1:00:00+02:30",
		"2026-10-13T1:00:00-02:30",
		"2026-10-13T1:00:00.123+02:30",
		"2026-10-13T1:00:00.123-02:30",
		"2026-10-13T12:00:00+24:00",
		"2026-10-13T12:00:00+00:60",
		"2026-10-13T12:00:00.1234567890Z",
		"2026-10-13T12:00:00,1234567890Z",
		"2026-10-13T12:00:00,123Z",
		"2026-10-13T12:00:00",
		"2026-10-13T12:00:00Z trailing",
		"infinity", "-infinity",
	} {
		t.Run(start, func(t *testing.T) {
			for _, endpoints := range [][2]string{{start, validEnd}, {"2026-10-12T12:00:00Z", start}} {
				for _, format := range []notation.Format{notation.ISO8601, notation.ISO80000, notation.Bourbaki} {
					input := endpoints[0] + "/" + endpoints[1]
					if format != notation.ISO8601 {
						input = "[" + endpoints[0] + "," + endpoints[1] + ")"
					}
					if format == notation.Bourbaki {
						input = "[" + endpoints[0] + "," + endpoints[1] + "["
					}
					period, err := notation.ParseInstant(input, format, temporal.Limits{})
					if !errors.Is(err, temporal.ErrParse) || !period.IsEmpty() {
						t.Errorf("unsupported endpoint admitted in %q: period=%+v err=%v", input, period, err)
					}
				}
			}
		})
	}
}

// FuzzInstantEndpointSelectedProfile checks accepted endpoints against an
// independent lexical profile and requires bounded, classified zero results on
// refusal. It does not claim complete RFC or cross-provider conformance.
func FuzzInstantEndpointSelectedProfile(f *testing.F) {
	for _, seed := range []string{
		"2026-10-13T1:00:00Z", "2026-10-13T12:00:00+24:00",
		"2026-10-13T1:00:00.123Z",
		"2026-10-13T1:00:00+02:30",
		"2026-10-13T1:00:00-02:30",
		"2026-10-13T1:00:00.123+02:30",
		"2026-10-13T1:00:00.123-02:30",
		"2026-10-13T12:00:00+00:60", "2026-10-13T12:00:00.1234567890Z",
		"2026-10-13T12:00:00,1234567890Z", "2026-10-13T12:00:00,123Z",
		"2026-10-13T12:00:00.123456789+02:30", "2026-10-13T12:00:00Z",
	} {
		f.Add(seed)
	}
	profile := regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}T([01][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9](\.[0-9]{1,9})?(Z|[+-]([01][0-9]|2[0-3]):[0-5][0-9])$`)
	f.Fuzz(func(t *testing.T, endpoint string) {
		if len(endpoint) > 4096 {
			t.Skip()
		}
		period, err := notation.ParseInstant(endpoint+"/9999-12-31T23:59:59Z", notation.ISO8601, temporal.Limits{})
		if err != nil {
			if !period.IsEmpty() || !errors.Is(err, temporal.ErrParse) || !utf8.ValidString(err.Error()) || len(err.Error()) > temporal.DefaultLimits().ErrorBytes {
				t.Fatalf("endpoint refusal lost diagnostic/zero-result contract: %+v, %v", period, err)
			}
			return
		}
		if !profile.MatchString(endpoint) {
			t.Fatalf("accepted endpoint outside the strict profile: %q", endpoint)
		}
		want, err := time.Parse(time.RFC3339Nano, endpoint)
		if err != nil || !period.Start().Equal(want) {
			t.Fatalf("accepted endpoint changed the represented instant: %v", err)
		}
	})
}

func TestInstantEndpointFractionLimitAndOffsetPositiveControls(t *testing.T) {
	for _, test := range []struct {
		start         string
		precision     int
		wantPrecision bool
	}{
		{"2026-10-13T12:00:00.123Z", 3, false},
		{"2026-10-13T01:00:00.123Z", 9, false},
		{"2026-10-13T01:00:00+02:30", 9, false},
		{"2026-10-13T01:00:00-02:30", 9, false},
		{"2026-10-13T12:00:00.1234Z", 3, true},
		{"2026-10-13T12:00:00,1234567890Z", 9, true},
		{"2026-10-13T12:00:00.123456789+02:30", 9, false},
		{"2026-10-13T12:00:00+23:59", 9, false},
		{"2026-10-13T12:00:00-23:59", 9, false},
		{"2026-10-13T12:00:00Z", 1, false},
	} {
		period, err := notation.ParseInstant(test.start+"/2026-10-15T12:00:00Z", notation.ISO8601, temporal.Limits{Precision: test.precision})
		if test.wantPrecision {
			if !errors.Is(err, temporal.ErrParse) || !errors.Is(err, temporal.ErrPrecision) || !period.IsEmpty() {
				t.Errorf("precision admission = %+v, %v", period, err)
			}
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
		want, err := time.Parse(time.RFC3339Nano, test.start)
		if err != nil {
			t.Fatal(err)
		}
		_, gotOffset := period.Start().Zone()
		_, wantOffset := want.Zone()
		if !period.Start().Equal(want) || gotOffset != wantOffset {
			t.Errorf("parsed endpoint lost instant or numeric offset")
		}
		encoded, err := notation.FormatInstant(period, notation.ISO8601, temporal.Limits{})
		if err != nil || encoded != test.start+"/2026-10-15T12:00:00Z" {
			t.Errorf("endpoint spelling = %q, %v", encoded, err)
		}
	}
}
