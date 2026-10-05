package dateperiod_test

import (
	"errors"
	"os"
	"testing"

	temporal "github.com/faustbrian/go-temporal/v2"
	"github.com/faustbrian/go-temporal/v2/dateperiod"
)

func TestDateSetUnionRefusesBeforeCombiningInputs(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "true" || os.Getenv("TEMPORAL_VERIFY_UNION_ADMISSION") != "true" {
		t.Skip("allocation admission proof requires explicit hosted execution")
	}

	first := mustDatePeriod(t, dayOffset(0), dayOffset(0), temporal.Closed)
	second := mustDatePeriod(t, dayOffset(2), dayOffset(2), temporal.Closed)
	periods := []dateperiod.Period{first, second}
	limits := temporal.Limits{InputPeriods: 1}
	left, err := dateperiod.NewSet(limits, first)
	if err != nil {
		t.Fatal(err)
	}
	right, err := dateperiod.NewSet(limits, second)
	if err != nil {
		t.Fatal(err)
	}
	refused, err := left.Union(right)
	if !errors.Is(err, temporal.ErrLimit) || refused.Len() != 0 {
		t.Fatalf("Union refusal = %v, %v; want empty result and ErrLimit", refused, err)
	}
	assertDateLimitError(t, err, "input_periods", 2, 1)
	refused, err = dateperiod.NewSet(limits, periods...)
	if !errors.Is(err, temporal.ErrLimit) || refused.Len() != 0 {
		t.Fatalf("NewSet refusal = %v, %v; want empty result and ErrLimit", refused, err)
	}
	assertDateLimitError(t, err, "input_periods", 2, 1)

	// Compare the same typed refusal, excluding all fixture construction.
	// Union must not allocate a combined slice before rejecting its count.
	var refusal error
	constructorAllocs := testing.AllocsPerRun(1, func() {
		_, refusal = dateperiod.NewSet(limits, periods...)
	})
	if !errors.Is(refusal, temporal.ErrLimit) {
		t.Fatalf("measured NewSet error = %v, want ErrLimit", refusal)
	}
	assertDateLimitError(t, refusal, "input_periods", 2, 1)
	unionAllocs := testing.AllocsPerRun(1, func() {
		_, refusal = left.Union(right)
	})
	if !errors.Is(refusal, temporal.ErrLimit) {
		t.Fatalf("measured Union error = %v, want ErrLimit", refusal)
	}
	assertDateLimitError(t, refusal, "input_periods", 2, 1)
	if unionAllocs > constructorAllocs {
		t.Fatalf("Union allocated before refusing input_periods: %g allocations, constructor refusal %g", unionAllocs, constructorAllocs)
	}
}
