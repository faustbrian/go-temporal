package timeofday_test

import (
	"errors"
	"testing"

	temporal "github.com/faustbrian/go-temporal/v2"
	"github.com/faustbrian/go-temporal/v2/timeofday"
)

func TestDailyUnionAdmitsNormalizedSegmentsBeforeCombining(t *testing.T) {
	limits := temporal.Limits{InputPeriods: 1}
	left, err := timeofday.NewIntervalSet(limits, mustInterval(t, 8, 9, temporal.ClosedOpen))
	if err != nil {
		t.Fatal(err)
	}
	right, err := timeofday.NewIntervalSet(temporal.Limits{}, mustInterval(t, 10, 11, temporal.ClosedOpen), mustInterval(t, 12, 13, temporal.ClosedOpen))
	if err != nil {
		t.Fatal(err)
	}
	beforeLeft, _ := timeofday.NewIntervalSet(limits, left.Intervals()...)
	beforeRight, _ := timeofday.NewIntervalSet(temporal.Limits{}, right.Intervals()...)
	result, err := left.Union(right)
	var limit *temporal.LimitError
	if !errors.Is(err, temporal.ErrLimit) || !errors.As(err, &limit) || limit.Field != "input_segments" || limit.Value != 3 || limit.Max != 2 || result.Len() != 0 {
		t.Fatalf("three segments under InputPeriods=1: result=%v error=%v; want zero result and typed 3/2 admission refusal", result.Intervals(), err)
	}
	if !left.Equal(beforeLeft) || !right.Equal(beforeRight) {
		t.Fatal("rejected union changed an operand")
	}
}

func TestDailyUnionInclusiveNormalizedSegmentAllowance(t *testing.T) {
	limits := temporal.Limits{InputPeriods: 1}
	left, err := timeofday.NewIntervalSet(limits, mustInterval(t, 8, 9, temporal.ClosedOpen))
	if err != nil {
		t.Fatal(err)
	}
	right, err := timeofday.NewIntervalSet(limits, mustInterval(t, 10, 11, temporal.ClosedOpen))
	if err != nil {
		t.Fatal(err)
	}
	union, err := left.Union(right)
	if err != nil || union.Len() != 2 || !union.Includes(hm(t, 8, 30)) || !union.Includes(hm(t, 10, 30)) || union.Includes(hm(t, 9, 30)) {
		t.Fatalf("inclusive same-budget union: %v, %v", union.Intervals(), err)
	}
	circular, err := timeofday.NewIntervalSet(limits, mustInterval(t, 22, 2, temporal.ClosedOpen))
	if err != nil {
		t.Fatal(err)
	}
	empty, err := timeofday.NewIntervalSet(limits)
	if err != nil {
		t.Fatal(err)
	}
	identity, err := circular.Union(empty)
	if err != nil || circular.Len() != 2 || !identity.Equal(circular) {
		t.Fatalf("circular expansion lost union identity: %v, %v", identity.Intervals(), err)
	}
}
