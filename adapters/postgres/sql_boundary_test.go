package temporalpostgres

import (
	"errors"
	"strings"
	"testing"

	temporal "github.com/faustbrian/go-temporal/v2"
)

func TestSQLTextExactByteLimitReachesRangeParsers(t *testing.T) {
	// The byte limit admits text independently of whether its range grammar is
	// valid. Reuse the same bounded payload for both supported SQL input types.
	input := strings.Repeat("x", 65536)
	for _, test := range []struct {
		name   string
		source any
	}{
		{"string", input},
		{"bytes", []byte(input)},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Run("admission", func(t *testing.T) {
				text, err := sqlText(test.source)
				if err != nil || text != input {
					t.Fatalf("exact-limit admission returned length %d, error %v", len(text), err)
				}
			})
			t.Run("instant receiver", func(t *testing.T) {
				var value InstantRange
				if err := value.Scan(`["2026-01-01T00:00:00Z","2026-01-02T00:00:00Z")`); err != nil {
					t.Fatal(err)
				}
				before := value
				if err := value.Scan(test.source); !errors.Is(err, temporal.ErrParse) || errors.Is(err, temporal.ErrLimit) {
					t.Fatalf("exact-limit Scan error = %v; want ErrParse, not ErrLimit", err)
				}
				if value != before {
					t.Fatal("failed Scan changed the instant receiver")
				}
			})
			t.Run("date receiver", func(t *testing.T) {
				var value DateRange
				if err := value.Scan("[2026-01-01,2026-01-02)"); err != nil {
					t.Fatal(err)
				}
				before := value
				if err := value.Scan(test.source); !errors.Is(err, temporal.ErrParse) || errors.Is(err, temporal.ErrLimit) {
					t.Fatalf("exact-limit Scan error = %v; want ErrParse, not ErrLimit", err)
				}
				if value != before {
					t.Fatal("failed Scan changed the date receiver")
				}
			})
		})
	}
}

func TestSQLTextAdmissionRejectsOverBudgetBeforeReturningText(t *testing.T) {
	// This checks admission/category only; pre-conversion ordering is established
	// by the owning sqlText source, not by an allocation assertion.
	input := make([]byte, temporal.DefaultLimits().ParseBytes+1)
	for _, source := range []any{input, string(input)} {
		text, err := sqlText(source)
		if !errors.Is(err, temporal.ErrLimit) || text != "" {
			t.Fatalf("sqlText = %q, %v; want empty text and ErrLimit", text, err)
		}
	}
}
