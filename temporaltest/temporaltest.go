// Package temporaltest provides canonical algebra fixtures and assertions for
// temporal package consumers. Assertions borrow TestingT and call its methods
// synchronously; nil implementations, panics and side effects remain caller-owned.
// Fixtures require representable caller-chosen base/unit offsets and do not
// independently certify all degenerate, calendar, timezone or compatibility cases.
package temporaltest

import (
	"time"

	temporal "github.com/faustbrian/go-temporal/v3"
	"github.com/faustbrian/go-temporal/v3/instant"
	"github.com/faustbrian/go-temporal/v3/timeofday"
)

// TestingT is the minimal testing contract used by assertions.
type TestingT interface {
	// Helper marks the calling assertion as a testing helper.
	Helper()
	// Errorf records a formatted failure; assertions do not require fatal behavior.
	Errorf(format string, arguments ...any)
}

// RelationCase is a canonical pair demonstrating one Allen relation.
type RelationCase struct {
	// Name is the canonical relation label for this fixture, not an independent
	// expected oracle.
	Name string
	// Left is the first half-open instant operand.
	Left instant.Period
	// Right is the second half-open instant operand.
	Right instant.Period
	// Want is the expected endpoint classification for this constructed
	// nondegenerate pair.
	Want temporal.Relation
}

// AllenCases builds one half-open nondegenerate fixture per Allen relation.
// unit must be positive; choose a unit/base whose small offset products and
// time.Time additions are representable. These deterministic cases are not an
// independent algebra oracle and do not exercise singleton or DST boundaries.
func AllenCases(base time.Time, unit time.Duration) ([]RelationCase, error) {
	if unit <= 0 {
		return nil, temporal.ErrStep
	}
	period := func(start, end int) instant.Period {
		value, _ := instant.Range(base.Add(time.Duration(start)*unit), base.Add(time.Duration(end)*unit))
		return value
	}
	definitions := [...]struct {
		relation   temporal.Relation
		leftStart  int
		leftEnd    int
		rightStart int
		rightEnd   int
	}{
		{temporal.Before, 0, 1, 2, 3},
		{temporal.Meets, 0, 1, 1, 2},
		{temporal.Overlaps, 0, 2, 1, 3},
		{temporal.Starts, 0, 1, 0, 2},
		{temporal.During, 1, 2, 0, 3},
		{temporal.Finishes, 1, 2, 0, 2},
		{temporal.Equal, 0, 2, 0, 2},
		{temporal.FinishedBy, 0, 2, 1, 2},
		{temporal.Contains, 0, 3, 1, 2},
		{temporal.StartedBy, 0, 2, 0, 1},
		{temporal.OverlappedBy, 1, 3, 0, 2},
		{temporal.MetBy, 1, 2, 0, 1},
		{temporal.After, 2, 3, 0, 1},
	}

	result := make([]RelationCase, len(definitions))
	for index, definition := range definitions {
		result[index] = RelationCase{
			Name:  definition.relation.String(),
			Left:  period(definition.leftStart, definition.leftEnd),
			Right: period(definition.rightStart, definition.rightEnd),
			Want:  definition.relation,
		}
	}
	return result, nil
}

// AssertRelation reports a mismatch or relation-domain error.
func AssertRelation(t TestingT, left, right instant.Period, want temporal.Relation) {
	t.Helper()
	got, err := left.RelationTo(right)
	if err != nil {
		t.Errorf("RelationTo() error = %v; want %s", err, want)
		return
	}
	if got != want {
		t.Errorf("RelationTo() = %s; want %s", got, want)
	}
}

// AssertInstantSetEqual compares normalized represented instant membership.
func AssertInstantSetEqual(t TestingT, got, want instant.Set) {
	t.Helper()
	if !got.Equal(want) {
		t.Errorf("instant sets differ: got %+v; want %+v", got.Periods(), want.Periods())
	}
}

// AssertDailySetEqual compares normalized represented daily membership.
func AssertDailySetEqual(t TestingT, got, want timeofday.IntervalSet) {
	t.Helper()
	if !got.Equal(want) {
		t.Errorf("daily sets differ: got %+v; want %+v", got.Intervals(), want.Intervals())
	}
}
