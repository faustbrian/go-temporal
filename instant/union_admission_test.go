package instant_test

import (
	"errors"
	"os"
	"testing"

	temporal "github.com/faustbrian/go-temporal/v2"
	"github.com/faustbrian/go-temporal/v2/instant"
)

func TestInstantSetUnionRefusesBeforeCombiningInputs(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "true" || os.Getenv("TEMPORAL_VERIFY_UNION_ADMISSION") != "true" {
		t.Skip("allocation admission proof requires explicit hosted execution")
	}

	// Four Period values exceed Go 1.27's 32-byte dynamic-make stack buffer:
	// each contains two time.Time values and Bounds (56 bytes on hosted 64-bit).
	// Recheck this witness when compiler thresholds or Period layout change.
	periods := []instant.Period{
		mustPeriod(t, 0, 0, temporal.Closed),
		mustPeriod(t, 2, 2, temporal.Closed),
		mustPeriod(t, 4, 4, temporal.Closed),
		mustPeriod(t, 6, 6, temporal.Closed),
	}
	limits := temporal.Limits{InputPeriods: 2}
	left, err := instant.NewSet(limits, periods[:2]...)
	if err != nil {
		t.Fatal(err)
	}
	right, err := instant.NewSet(limits, periods[2:]...)
	if err != nil {
		t.Fatal(err)
	}
	assertRefusal := func(result instant.Set, err error) {
		t.Helper()
		var limitError *temporal.LimitError
		if result.Len() != 0 || !errors.Is(err, temporal.ErrLimit) || !errors.As(err, &limitError) {
			t.Fatalf("refusal = %v, %v; want empty result and LimitError", result, err)
		}
		if limitError.Field != "input_periods" || limitError.Value != 4 || limitError.Max != 2 {
			t.Fatalf("LimitError = %+v; want input_periods=4, maximum 2", limitError)
		}
	}
	assertRefusal(left.Union(right))
	assertRefusal(instant.NewSet(limits, periods...))

	// Compare the same typed refusal, excluding all fixture construction.
	// Union must not allocate a combined slice before rejecting its count.
	var refusal error
	constructorAllocs := testing.AllocsPerRun(1, func() {
		_, refusal = instant.NewSet(limits, periods...)
	})
	assertRefusal(instant.Set{}, refusal)
	unionAllocs := testing.AllocsPerRun(1, func() {
		_, refusal = left.Union(right)
	})
	assertRefusal(instant.Set{}, refusal)
	if unionAllocs > constructorAllocs {
		t.Fatalf("Union allocated before refusing input_periods: %g allocations, constructor refusal %g", unionAllocs, constructorAllocs)
	}
}
