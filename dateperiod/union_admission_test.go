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

	// Four 10-byte Period values exceed Go 1.27's 32-byte dynamic-make
	// stack buffer; two values can hide the unnecessary copy from heap counts.
	// Recheck this witness when compiler thresholds or Period layout change.
	periods := []dateperiod.Period{
		mustDatePeriod(t, dayOffset(0), dayOffset(0), temporal.Closed),
		mustDatePeriod(t, dayOffset(2), dayOffset(2), temporal.Closed),
		mustDatePeriod(t, dayOffset(4), dayOffset(4), temporal.Closed),
		mustDatePeriod(t, dayOffset(6), dayOffset(6), temporal.Closed),
	}
	limits := temporal.Limits{InputPeriods: 2}
	left, err := dateperiod.NewSet(limits, periods[:2]...)
	if err != nil {
		t.Fatal(err)
	}
	right, err := dateperiod.NewSet(limits, periods[2:]...)
	if err != nil {
		t.Fatal(err)
	}
	refused, err := left.Union(right)
	if !errors.Is(err, temporal.ErrLimit) || refused.Len() != 0 {
		t.Fatalf("Union refusal = %v, %v; want empty result and ErrLimit", refused, err)
	}
	assertDateLimitError(t, err, "input_periods", 4, 2)
	refused, err = dateperiod.NewSet(limits, periods...)
	if !errors.Is(err, temporal.ErrLimit) || refused.Len() != 0 {
		t.Fatalf("NewSet refusal = %v, %v; want empty result and ErrLimit", refused, err)
	}
	assertDateLimitError(t, err, "input_periods", 4, 2)

	// Compare the same typed refusal, excluding all fixture construction.
	// Union must not allocate a combined slice before rejecting its count.
	var refusal error
	constructorAllocs := testing.AllocsPerRun(1, func() {
		_, refusal = dateperiod.NewSet(limits, periods...)
	})
	if !errors.Is(refusal, temporal.ErrLimit) {
		t.Fatalf("measured NewSet error = %v, want ErrLimit", refusal)
	}
	assertDateLimitError(t, refusal, "input_periods", 4, 2)
	unionAllocs := testing.AllocsPerRun(1, func() {
		_, refusal = left.Union(right)
	})
	if !errors.Is(refusal, temporal.ErrLimit) {
		t.Fatalf("measured Union error = %v, want ErrLimit", refusal)
	}
	assertDateLimitError(t, refusal, "input_periods", 4, 2)
	if unionAllocs > constructorAllocs {
		t.Fatalf("Union allocated before refusing input_periods: %g allocations, constructor refusal %g", unionAllocs, constructorAllocs)
	}
}
