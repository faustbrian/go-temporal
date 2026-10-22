package notation_test

import (
	"testing"

	temporal "github.com/faustbrian/go-temporal/v3"
	"github.com/faustbrian/go-temporal/v3/notation"
	"github.com/faustbrian/go-temporal/v3/timeofday"
)

// FuzzInstantNotation exercises interval text and selected notation alphabets.
// Successful parses must format and reparse with the same represented set and bounds.
func FuzzInstantNotation(f *testing.F) {
	for _, seed := range []struct {
		value  string
		format uint8
	}{
		{"2026-01-02T03:04:05Z/2026-01-03T03:04:05Z", uint8(notation.ISO8601)},
		{"[2026-01-02T03:04:05Z,2026-01-03T03:04:05Z)", uint8(notation.ISO80000)},
		{"]2026-01-02T03:04:05.123456789Z,2026-01-03T03:04:05Z]", uint8(notation.Bourbaki)},
		{"bad", 255},
	} {
		f.Add(seed.value, seed.format)
	}
	f.Fuzz(func(t *testing.T, input string, formatValue uint8) {
		format := fuzzFormat(formatValue)
		period, err := notation.ParseInstant(input, format, temporal.Limits{})
		if err != nil {
			return
		}
		encoded, err := notation.FormatInstant(period, format, temporal.Limits{})
		if err != nil {
			t.Fatalf("FormatInstant(): %v", err)
		}
		roundTrip, err := notation.ParseInstant(encoded, format, temporal.Limits{})
		if err != nil || !roundTrip.SetEqual(period) || roundTrip.Bounds() != period.Bounds() {
			t.Fatalf("round trip %q = %+v, %v", encoded, roundTrip, err)
		}
	})
}

// FuzzDailyIntervalNotation exercises local interval text including explicit kinds.
// Successful parses must preserve structural interval identity through formatting.
func FuzzDailyIntervalNotation(f *testing.F) {
	for _, seed := range []struct {
		value  string
		format uint8
	}{
		{"22:00/02:30", uint8(notation.ISO8601)},
		{"(08:00:00.1,17:00:00]", uint8(notation.ISO80000)},
		{"[00:00:00,24:00]", uint8(notation.ISO80000)},
		{"]08:00,17:00[", uint8(notation.Bourbaki)},
		{"bad", 255},
	} {
		f.Add(seed.value, seed.format)
	}
	f.Fuzz(func(t *testing.T, input string, formatValue uint8) {
		format := fuzzFormat(formatValue)
		interval, err := notation.ParseDailyInterval(input, format, temporal.Limits{})
		if err != nil {
			return
		}
		encoded, err := notation.FormatDailyInterval(interval, format, temporal.Limits{})
		if err != nil {
			t.Fatalf("FormatDailyInterval(): %v", err)
		}
		roundTrip, err := notation.ParseDailyInterval(encoded, format, temporal.Limits{})
		if err != nil || !roundTrip.Equal(interval) {
			t.Fatalf("round trip %q = %+v, %v", encoded, roundTrip, err)
		}
	})
}

// FuzzDateNotation exercises civil interval text and notation alphabets.
// Successful parses must preserve represented dates and bounds through formatting.
func FuzzDateNotation(f *testing.F) {
	for _, seed := range []struct {
		value  string
		format uint8
	}{
		{"2026-01-01/2027-01-01", uint8(notation.ISO8601)},
		{"[2026-01-01,2026-12-31]", uint8(notation.ISO80000)},
		{"]2026-01-01,2026-12-31[", uint8(notation.Bourbaki)},
		{"bad", 255},
	} {
		f.Add(seed.value, seed.format)
	}
	f.Fuzz(func(t *testing.T, input string, formatValue uint8) {
		format := fuzzFormat(formatValue)
		period, err := notation.ParseDate(input, format, temporal.Limits{})
		if err != nil {
			return
		}
		encoded, err := notation.FormatDate(period, format, temporal.Limits{})
		if err != nil {
			t.Fatalf("FormatDate(): %v", err)
		}
		roundTrip, err := notation.ParseDate(encoded, format, temporal.Limits{})
		if err != nil || !roundTrip.SetEqual(period) || roundTrip.Bounds() != period.Bounds() {
			t.Fatalf("round trip %q = %+v, %v", encoded, roundTrip, err)
		}
	})
}

func fuzzFormat(value uint8) notation.Format {
	formats := [...]notation.Format{notation.ISO8601, notation.ISO80000, notation.Bourbaki}
	return formats[int(value)%len(formats)]
}

// FuzzFixedDurationNotation exercises fixed elapsed ISO-subset text.
// Admitted values must retain exact duration through formatting and reparsing.
func FuzzFixedDurationNotation(f *testing.F) {
	for _, seed := range []string{"PT0S", "-P1DT2H3M4.123456789S", "P1M", "bad"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		duration, err := notation.ParseDuration(input, temporal.Limits{})
		if err != nil {
			return
		}
		encoded, err := notation.FormatDuration(duration, temporal.Limits{})
		if err != nil {
			t.Fatalf("FormatDuration(): %v", err)
		}
		roundTrip, err := notation.ParseDuration(encoded, temporal.Limits{})
		if err != nil || roundTrip != duration {
			t.Fatalf("round trip %q = %v, %v", encoded, roundTrip, err)
		}
	})
}

// FuzzLocalTime exercises local-time text, including malformed seed inputs.
// Admitted values must retain their String presentation on reparsing.
func FuzzLocalTime(f *testing.F) {
	for _, seed := range []string{"00:00", "24:00", "12:34:56.123456789", "bad"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		value, err := timeofday.Parse(input, temporal.Limits{})
		if err != nil {
			return
		}
		roundTrip, err := timeofday.Parse(value.String(), temporal.Limits{})
		if err != nil || roundTrip.String() != value.String() {
			t.Fatalf("round trip %q = %v, %v", value, roundTrip, err)
		}
	})
}

// BenchmarkInstantNotationParsing measures parsing one fixed fractional bracketed interval.
// Input construction is excluded; allocations and timing compare like workloads.
func BenchmarkInstantNotationParsing(b *testing.B) {
	input := "[2026-01-02T03:04:05.123456789Z,2026-01-03T04:05:06.987654321Z)"
	b.ReportAllocs()
	for b.Loop() {
		if _, err := notation.ParseInstant(input, notation.ISO80000, temporal.Limits{}); err != nil {
			b.Fatal(err)
		}
	}
}
