package temporal_test

import (
	"testing"
	"time"

	calendar "github.com/faustbrian/go-calendar/v2"
	temporal "github.com/faustbrian/go-temporal/v2"
	"github.com/faustbrian/go-temporal/v2/dateperiod"
	"github.com/faustbrian/go-temporal/v2/instant"
)

func TestCoincidentClosedSingletonRelationAndSetIdentity(t *testing.T) {
	at := time.Date(2026, time.October, 13, 12, 0, 0, 0, time.UTC)
	left := instant.Point(at)
	right := instant.Point(at.In(time.FixedZone("same instant", 3600)))
	for _, pair := range [][2]instant.Period{{left, right}, {right, left}} {
		relation, err := pair[0].RelationTo(pair[1])
		if err != nil || relation != temporal.Equal || relation.Converse() != temporal.Equal {
			t.Errorf("coincident instant relation = %v, %v; want self-converse Equal", relation, err)
		}
		intersection, ok := pair[0].Intersect(pair[1])
		if !ok || !intersection.SetEqual(pair[0]) || !pair[0].Contains(pair[1]) || len(pair[0].Subtract(pair[1])) != 0 {
			t.Errorf("coincident singleton lost exact set identity")
		}
		for _, offset := range []time.Duration{-time.Nanosecond, 0, time.Nanosecond} {
			if got := intersection.Includes(at.Add(offset)); got != (offset == 0) {
				t.Errorf("intersection membership at offset %v = %v", offset, got)
			}
		}
	}
	day := calendar.MustDate(2026, time.October, 13)
	date, err := dateperiod.New(day, day, temporal.Closed)
	if err != nil {
		t.Fatal(err)
	}
	relation, err := date.RelationTo(date)
	if err != nil || relation != temporal.Equal || relation.Converse() != temporal.Equal {
		t.Errorf("coincident civil singleton relation = %v, %v; want Equal", relation, err)
	}
	intersection, ok := date.Intersect(date)
	if !ok || !intersection.SetEqual(date) || len(date.Subtract(date)) != 0 {
		t.Error("civil singleton set identity failed")
	}
	for _, offset := range []int{-1, 0, 1} {
		probe, err := day.AddDays(offset)
		if err != nil {
			t.Fatal(err)
		}
		if got := intersection.Includes(probe); got != (offset == 0) {
			t.Errorf("civil membership at %d = %v", offset, got)
		}
	}
}

func TestDenseIntervalAlgebraMatchesIndependentFiniteMembershipOracle(t *testing.T) {
	origin := time.Date(2026, time.October, 13, 12, 0, 0, 0, time.UTC)
	type interval struct {
		lower, upper int
		bounds       temporal.Bounds
	}
	member := func(v interval, n int) bool {
		return (n > v.lower || n == v.lower && v.bounds.IncludesStart()) &&
			(n < v.upper || n == v.upper && v.bounds.IncludesEnd())
	}
	for _, bounds := range temporal.AllBounds() {
		for _, otherBounds := range temporal.AllBounds() {
			for _, endpoints := range [][4]int{{0, 0, 0, 0}, {0, 2, 2, 4}, {0, 4, 1, 3}, {0, 2, 1, 4}, {0, 1, 3, 4}} {
				a, b := interval{endpoints[0], endpoints[1], bounds}, interval{endpoints[2], endpoints[3], otherBounds}
				left, err := instant.New(origin.Add(time.Duration(a.lower)*time.Nanosecond), origin.Add(time.Duration(a.upper)*time.Nanosecond), a.bounds)
				if err != nil {
					t.Fatal(err)
				}
				right, err := instant.New(origin.Add(time.Duration(b.lower)*time.Nanosecond), origin.Add(time.Duration(b.upper)*time.Nanosecond), b.bounds)
				if err != nil {
					t.Fatal(err)
				}
				intersection, ok := left.Intersect(right)
				difference := left.Subtract(right)
				union, err := left.Union(right, temporal.Limits{})
				if err != nil {
					t.Fatal(err)
				}
				for n := -1; n <= 5; n++ {
					probe := origin.Add(time.Duration(n) * time.Nanosecond)
					inDifference := false
					for _, fragment := range difference {
						inDifference = inDifference || fragment.Includes(probe)
					}
					if (ok && intersection.Includes(probe)) != (member(a, n) && member(b, n)) ||
						inDifference != (member(a, n) && !member(b, n)) || union.Includes(probe) != (member(a, n) || member(b, n)) {
						t.Fatalf("set membership differs at %d for %+v, %+v", n, a, b)
					}
				}
			}
		}
	}
}
