package temporal_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	temporal "github.com/faustbrian/go-temporal"
	newconfig "github.com/faustbrian/go-temporal/adapters/config"
	newpostgres "github.com/faustbrian/go-temporal/adapters/postgres"
	newvalidation "github.com/faustbrian/go-temporal/adapters/validation"
	newwire "github.com/faustbrian/go-temporal/adapters/wire"
	"github.com/faustbrian/go-temporal/instant"
	oldpostgres "github.com/faustbrian/go-temporal/postgres"
	oldconfig "github.com/faustbrian/go-temporal/temporalconfig"
	oldvalidation "github.com/faustbrian/go-temporal/temporalvalidation"
	oldwire "github.com/faustbrian/go-temporal/temporalwire"
	validation "github.com/faustbrian/go-validation"
)

func TestSuccessorNamedTypesOwnTheirPackageIdentity(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		value any
		path  string
	}{
		"config instant":   {newconfig.InstantPeriod{}, "github.com/faustbrian/go-temporal/adapters/config"},
		"config date":      {newconfig.DatePeriod{}, "github.com/faustbrian/go-temporal/adapters/config"},
		"config daily":     {newconfig.DailyInterval{}, "github.com/faustbrian/go-temporal/adapters/config"},
		"config time":      {newconfig.Time{}, "github.com/faustbrian/go-temporal/adapters/config"},
		"config duration":  {newconfig.Duration{}, "github.com/faustbrian/go-temporal/adapters/config"},
		"wire kind":        {newwire.Kind(""), "github.com/faustbrian/go-temporal/adapters/wire"},
		"wire document":    {newwire.Document{}, "github.com/faustbrian/go-temporal/adapters/wire"},
		"wire collection":  {newwire.CollectionDocument{}, "github.com/faustbrian/go-temporal/adapters/wire"},
		"postgres instant": {newpostgres.InstantRange{}, "github.com/faustbrian/go-temporal/adapters/postgres"},
		"postgres date":    {newpostgres.DateRange{}, "github.com/faustbrian/go-temporal/adapters/postgres"},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if got := reflect.TypeOf(test.value).PkgPath(); got != test.path {
				t.Fatalf("PkgPath = %q, want %q", got, test.path)
			}
		})
	}

	if reflect.TypeOf(oldconfig.InstantPeriod{}).PkgPath() == reflect.TypeOf(newconfig.InstantPeriod{}).PkgPath() ||
		reflect.TypeOf(oldwire.Document{}).PkgPath() == reflect.TypeOf(newwire.Document{}).PkgPath() ||
		reflect.TypeOf(oldpostgres.InstantRange{}).PkgPath() == reflect.TypeOf(newpostgres.InstantRange{}).PkgPath() {
		t.Fatal("legacy and successor named types must remain distinct")
	}
}

func TestLegacyAndSuccessorAdaptersPreserveBehavior(t *testing.T) {
	t.Parallel()

	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	period, err := instant.Range(start, start.Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}

	oldDocument, err := oldwire.FromInstant(period, temporal.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	newDocument, err := newwire.FromInstant(period, temporal.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	oldBytes, err := oldwire.Marshal(oldDocument, temporal.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	newBytes, err := newwire.Marshal(newDocument, temporal.Limits{})
	if err != nil || string(oldBytes) != string(newBytes) {
		t.Fatalf("wire old=%s new=%s err=%v", oldBytes, newBytes, err)
	}

	oldRange, err := oldpostgres.NewInstantRange(period)
	if err != nil {
		t.Fatal(err)
	}
	newRange, err := newpostgres.NewInstantRange(period)
	if err != nil {
		t.Fatal(err)
	}
	oldValue, oldErr := oldRange.Value()
	newValue, newErr := newRange.Value()
	if oldErr != nil || newErr != nil || oldValue != newValue {
		t.Fatalf("postgres old=%v/%v new=%v/%v", oldValue, oldErr, newValue, newErr)
	}

	oldText, err := oldconfig.NewInstantPeriod(period).MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	newText, err := newconfig.NewInstantPeriod(period).MarshalText()
	if err != nil || string(oldText) != string(newText) {
		t.Fatalf("config old=%s new=%s err=%v", oldText, newText, err)
	}

	vctx, err := validation.NewContext(validation.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	oldReport := oldvalidation.InstantNonEmpty().Validate(vctx, instant.Period{})
	newReport := newvalidation.InstantNonEmpty().Validate(vctx, instant.Period{})
	for name, report := range map[string]validation.Report{"old": oldReport, "new": newReport} {
		merged := validation.ContextReport(vctx, ctx).Merge(report)
		if !errors.Is(merged.Err(), context.Canceled) || !errors.Is(merged.Err(), validation.ErrInvalid) ||
			!merged.HasCode("temporal_empty") {
			t.Fatalf("%s merged report = %v, err=%v", name, merged, merged.Err())
		}
	}
}
