package timeofday

import (
	"time"

	calendar "github.com/faustbrian/go-calendar/v2"
	calendartz "github.com/faustbrian/go-calendar/v2/timezone"
	temporal "github.com/faustbrian/go-temporal/v2"
	"github.com/faustbrian/go-temporal/v2/instant"
)

// ToInstant resolves a daily interval beginning on date. Circular intervals
// end on the following civil date. Full-day conversion always uses the next
// civil midnight as an exclusive end boundary.
func (i Interval) ToInstant(date calendar.Date, location *time.Location, policy calendartz.Resolution) (instant.Period, error) {
	if i.kind == FullDayKind {
		start, end, err := calendartz.DayRange(date, location, policy)
		if err != nil {
			return instant.Period{}, err
		}
		return instant.Range(start, end)
	}
	start, err := i.start.Apply(date, location, policy)
	if err != nil {
		return instant.Period{}, err
	}
	endDate := date
	if i.kind == Circular {
		endDate, err = date.AddDays(1)
		if err != nil {
			return instant.Period{}, err
		}
	}
	end, err := i.end.Apply(endDate, location, policy)
	if err != nil {
		return instant.Period{}, err
	}
	bounds := i.bounds
	if i.kind == CollapsedKind {
		bounds = temporal.Open
	}
	return instant.New(start, end, bounds)
}
