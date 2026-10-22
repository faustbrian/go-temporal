// Package instant provides immutable dense intervals over time.Time instants.
// New validates ordering/bounds and strips monotonic readings, retaining endpoint
// locations. The zero Period is anchored empty; Point creates a closed singleton.
// Endpoint Allen classification is separate from represented-set membership.
//
// Set copies input/output slices and normalizes endpoints; zero sets are empty
// with default limits. Operations return new values and use receiver limits.
// Iterators/callbacks execute synchronously, own no goroutines, and cannot preempt
// blocking or panicking caller code. Scalar algebra owns no clock or persistence;
// explicit civil snapping delegates to the pinned Calendar v2 timezone contract.
package instant

import (
	"fmt"
	"time"

	temporal "github.com/faustbrian/go-temporal/v3"
)

// Period is an immutable bounded interval between two instants. Construction
// strips process-local monotonic readings while retaining locations.
type Period struct {
	start  time.Time
	end    time.Time
	bounds temporal.Bounds
}

// New constructs a bounded period. Equal endpoints represent a singleton only
// when both endpoints are included; other equal-endpoint modes are empty.
func New(start, end time.Time, bounds temporal.Bounds) (Period, error) {
	if !bounds.Valid() {
		return Period{}, fmt.Errorf("%w: %d", temporal.ErrBounds, bounds)
	}

	start = stripMonotonic(start)
	end = stripMonotonic(end)
	if end.Before(start) {
		return Period{}, fmt.Errorf("%w: end precedes start", temporal.ErrReversed)
	}

	return Period{start: start, end: end, bounds: bounds}, nil
}

// Range constructs a canonical start-inclusive, end-exclusive period.
func Range(start, end time.Time) (Period, error) {
	return New(start, end, temporal.ClosedOpen)
}

// Start returns the start instant without a monotonic reading.
func (p Period) Start() time.Time {
	return p.start
}

// End returns the end instant without a monotonic reading.
func (p Period) End() time.Time {
	return p.end
}

// Bounds returns the endpoint inclusion mode.
func (p Period) Bounds() temporal.Bounds {
	return p.bounds
}

// IsEmpty reports whether the period represents the empty set.
func (p Period) IsEmpty() bool {
	return p.start.Equal(p.end) && p.bounds != temporal.Closed
}

// IsSingleton reports whether the period contains exactly one instant.
func (p Period) IsSingleton() bool {
	return p.start.Equal(p.end) && p.bounds == temporal.Closed
}

// Includes reports whether the instant is a member of the period.
func (p Period) Includes(value time.Time) bool {
	if p.IsEmpty() || value.Before(p.start) || value.After(p.end) {
		return false
	}
	if value.Equal(p.start) && !p.bounds.IncludesStart() {
		return false
	}
	if value.Equal(p.end) && !p.bounds.IncludesEnd() {
		return false
	}

	return true
}

// RelationTo classifies declared endpoint ordering; empty operands return
// RelationInvalid/ErrEmpty. Inclusion is not used beyond emptiness. Matching
// endpoints are Equal, including coincident closed singletons. Represented
// membership is separate from endpoint classification and structural equality.
func (p Period) RelationTo(other Period) (temporal.Relation, error) {
	if p.IsEmpty() || other.IsEmpty() {
		return temporal.RelationInvalid, temporal.ErrEmpty
	}
	if p.start.Equal(other.start) && p.end.Equal(other.end) {
		return temporal.Equal, nil
	}

	if p.end.Before(other.start) {
		return temporal.Before, nil
	}
	if p.end.Equal(other.start) {
		return temporal.Meets, nil
	}
	if p.start.After(other.end) {
		return temporal.After, nil
	}
	if p.start.Equal(other.end) {
		return temporal.MetBy, nil
	}

	startComparison := p.start.Compare(other.start)
	endComparison := p.end.Compare(other.end)

	switch startComparison {
	case -1:
		switch endComparison {
		case -1:
			return temporal.Overlaps, nil
		case 0:
			return temporal.FinishedBy, nil
		default:
			return temporal.Contains, nil
		}
	case 0:
		switch endComparison {
		case -1:
			return temporal.Starts, nil
		default:
			return temporal.StartedBy, nil
		}
	default:
		switch endComparison {
		case -1:
			return temporal.During, nil
		case 0:
			return temporal.Finishes, nil
		default:
			return temporal.OverlappedBy, nil
		}
	}
}

// Abuts reports whether one end is equal to the other start, regardless of
// whether the shared endpoint is included.
func (p Period) Abuts(other Period) bool {
	return p.end.Equal(other.start) || other.end.Equal(p.start)
}

// Borders reports whether abutting periods both include their shared endpoint.
func (p Period) Borders(other Period) bool {
	if p.end.Equal(other.start) {
		return p.bounds.IncludesEnd() && other.bounds.IncludesStart()
	}
	if other.end.Equal(p.start) {
		return other.bounds.IncludesEnd() && p.bounds.IncludesStart()
	}

	return false
}

// Meets reports whether periods abut without sharing the adjacent endpoint.
func (p Period) Meets(other Period) bool {
	return p.Abuts(other) && !p.Borders(other)
}

func stripMonotonic(value time.Time) time.Time {
	return value.Round(0)
}
