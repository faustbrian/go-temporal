package instant

import (
	"iter"
	"sort"
	"time"

	temporal "github.com/faustbrian/go-temporal/v2"
)

// All returns a repeatable iterator over stable normalized order. It calls yield
// synchronously until false; callback panics or blocking propagate to the caller.
func (s Set) All() iter.Seq[Period] {
	return func(yield func(Period) bool) {
		for _, period := range s.periods {
			if !yield(period) {
				return
			}
		}
	}
}

// Search returns the containing period index or the stable insertion index.
func (s Set) Search(value time.Time) (int, bool) {
	index := sort.Search(len(s.periods), func(index int) bool {
		return !s.periods[index].end.Before(value)
	})
	if index < len(s.periods) && s.periods[index].Includes(value) {
		return index, true
	}
	return index, false
}

// Transform calls mapper synchronously in stable order and normalizes with the
// receiver limits. A nil mapper fails with ErrUnsupported; mapper or normalization
// failure returns a zero set. Callback side effects are not rolled back.
func (s Set) Transform(mapper func(Period) (Period, error)) (Set, error) {
	if mapper == nil {
		return Set{}, temporal.ErrUnsupported
	}
	periods := make([]Period, len(s.periods))
	for index, period := range s.periods {
		mapped, err := mapper(period)
		if err != nil {
			return Set{}, err
		}
		periods[index] = mapped
	}
	return NewSet(s.effectiveLimits(), periods...)
}

// Reduce calls reducer synchronously in stable order. Nil reducers and callback
// errors return initial with an error, not a partial accumulator. T is copied
// only as a Go value: mutations through aliases or callback effects are not undone.
func Reduce[T any](set Set, initial T, reducer func(T, Period) (T, error)) (T, error) {
	if reducer == nil {
		return initial, temporal.ErrUnsupported
	}
	result := initial
	for _, period := range set.periods {
		var err error
		result, err = reducer(result, period)
		if err != nil {
			return initial, err
		}
	}
	return result, nil
}
