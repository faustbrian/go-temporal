package temporalvalidation_test

import (
	"errors"
	"testing"
	"time"

	calendar "github.com/faustbrian/go-calendar/v2"
	temporal "github.com/faustbrian/go-temporal/v2"
	temporalvalidation "github.com/faustbrian/go-temporal/v2/adapters/validation"
	"github.com/faustbrian/go-temporal/v2/dateperiod"
	"github.com/faustbrian/go-temporal/v2/instant"
	"github.com/faustbrian/go-temporal/v2/timeofday"
	validation "github.com/faustbrian/go-validation/v2"
)

func TestNonEmptyValidatorsReturnStableViolations(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	emptyInstant, _ := instant.New(now, now, temporal.Open)
	fullInstant, _ := instant.New(now, now, temporal.Closed)
	date := calendar.MustDate(2026, time.January, 2)
	emptyDate, _ := dateperiod.New(date, date, temporal.Open)
	fullDate, _ := dateperiod.New(date, date, temporal.Closed)
	anchor, _ := timeofday.Parse("08:00", temporal.Limits{})
	emptyDaily := timeofday.Collapsed(anchor)
	fullDaily := timeofday.FullDay()

	tests := []struct {
		name string
		pass validation.Report
		fail validation.Report
	}{
		{"instant", temporalvalidation.InstantNonEmpty().Validate(validation.Context{}, fullInstant), temporalvalidation.InstantNonEmpty().Validate(validation.Context{}, emptyInstant)},
		{"date", temporalvalidation.DateNonEmpty().Validate(validation.Context{}, fullDate), temporalvalidation.DateNonEmpty().Validate(validation.Context{}, emptyDate)},
		{"daily", temporalvalidation.DailyNonEmpty().Validate(validation.Context{}, fullDaily), temporalvalidation.DailyNonEmpty().Validate(validation.Context{}, emptyDaily)},
	}
	for _, test := range tests {
		if !test.pass.Empty() || !test.fail.HasErrors() || !test.fail.HasCode("temporal_empty") {
			t.Fatalf("%s reports = pass:%v fail:%v", test.name, test.pass, test.fail)
		}
	}
}

func TestRangeValidatorsUseSemanticTemporalOrdering(t *testing.T) {
	t.Parallel()

	minimum, _ := timeofday.Parse("08:00", temporal.Limits{})
	maximum, _ := timeofday.Parse("17:00", temporal.Limits{})
	inside, _ := timeofday.Parse("12:00", temporal.Limits{})
	out, _ := timeofday.Parse("18:00", temporal.Limits{})
	timeRule, err := temporalvalidation.TimeBetween(minimum, maximum)
	if err != nil {
		t.Fatalf("TimeBetween(): %v", err)
	}
	if !timeRule.Validate(validation.Context{}, inside).Empty() || !timeRule.Validate(validation.Context{}, out).HasCode("time_of_day_range") {
		t.Fatal("time range validator accepted or rejected the wrong value")
	}
	if _, err := temporalvalidation.TimeBetween(maximum, minimum); !errors.Is(err, temporal.ErrReversed) {
		t.Fatalf("TimeBetween(reversed) error = %v, want ErrReversed", err)
	}

	durationRule, err := temporalvalidation.DurationBetween(
		timeofday.NewDuration(time.Minute),
		timeofday.NewDuration(time.Hour),
	)
	if err != nil {
		t.Fatalf("DurationBetween(): %v", err)
	}
	if !durationRule.Validate(validation.Context{}, timeofday.NewDuration(30*time.Minute)).Empty() ||
		!durationRule.Validate(validation.Context{}, timeofday.NewDuration(2*time.Hour)).HasCode("fixed_duration_range") {
		t.Fatal("duration range validator accepted or rejected the wrong value")
	}
	if _, err := temporalvalidation.DurationBetween(timeofday.NewDuration(time.Hour), timeofday.NewDuration(time.Minute)); !errors.Is(err, temporal.ErrReversed) {
		t.Fatalf("DurationBetween(reversed) error = %v, want ErrReversed", err)
	}
}

func TestRangeValidatorsIncludeBothExactEndpointsAndSingletonRanges(t *testing.T) {
	t.Parallel()

	minimum, _ := timeofday.Parse("08:00", temporal.Limits{})
	maximum, _ := timeofday.Parse("17:00", temporal.Limits{})
	timeRule, err := temporalvalidation.TimeBetween(minimum, maximum)
	if err != nil {
		t.Fatalf("TimeBetween(): %v", err)
	}
	for _, endpoint := range []timeofday.Time{minimum, maximum} {
		if report := timeRule.Validate(validation.Context{}, endpoint); !report.Empty() {
			t.Fatalf("TimeBetween rejected endpoint %v: %v", endpoint, report)
		}
	}
	singleTime, err := temporalvalidation.TimeBetween(minimum, minimum)
	if err != nil {
		t.Fatalf("TimeBetween(singleton): %v", err)
	}
	if !singleTime.Validate(validation.Context{}, minimum).Empty() ||
		!singleTime.Validate(validation.Context{}, maximum).HasCode("time_of_day_range") {
		t.Fatal("singleton time range did not accept only its endpoint")
	}

	minimumDuration := timeofday.NewDuration(time.Minute)
	maximumDuration := timeofday.NewDuration(time.Hour)
	durationRule, err := temporalvalidation.DurationBetween(minimumDuration, maximumDuration)
	if err != nil {
		t.Fatalf("DurationBetween(): %v", err)
	}
	for _, endpoint := range []timeofday.Duration{minimumDuration, maximumDuration} {
		if report := durationRule.Validate(validation.Context{}, endpoint); !report.Empty() {
			t.Fatalf("DurationBetween rejected endpoint %v: %v", endpoint, report)
		}
	}
	singleDuration, err := temporalvalidation.DurationBetween(minimumDuration, minimumDuration)
	if err != nil {
		t.Fatalf("DurationBetween(singleton): %v", err)
	}
	if !singleDuration.Validate(validation.Context{}, minimumDuration).Empty() ||
		!singleDuration.Validate(validation.Context{}, maximumDuration).HasCode("fixed_duration_range") {
		t.Fatal("singleton duration range did not accept only its endpoint")
	}
}
