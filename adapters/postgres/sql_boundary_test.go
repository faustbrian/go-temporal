package temporalpostgres

import (
	"errors"
	"testing"

	temporal "github.com/faustbrian/go-temporal/v2"
)

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
