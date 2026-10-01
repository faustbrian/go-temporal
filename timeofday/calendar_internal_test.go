package timeofday

import (
	"testing"
	"time"

	calendar "github.com/faustbrian/go-calendar/v2"
	calendartz "github.com/faustbrian/go-calendar/v2/timezone"
)

func TestApplyRejectsInternallyInvalidTimeValue(t *testing.T) {
	invalid := Time{offset: 25 * time.Hour}
	if _, err := invalid.Apply(calendar.MustDate(2026, time.January, 1), time.UTC, calendartz.Reject); err == nil {
		t.Fatal("Apply accepted invalid internal time")
	}
}
