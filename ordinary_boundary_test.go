package temporal_test

import (
	"database/sql/driver"
	"errors"
	"strings"
	"testing"
	"time"

	calendar "github.com/faustbrian/go-calendar/v2"
	temporal "github.com/faustbrian/go-temporal/v2"
	postgres "github.com/faustbrian/go-temporal/v2/adapters/postgres"
	wire "github.com/faustbrian/go-temporal/v2/adapters/wire"
	"github.com/faustbrian/go-temporal/v2/dateperiod"
	"github.com/faustbrian/go-temporal/v2/instant"
	"github.com/faustbrian/go-temporal/v2/timeofday"
)

func TestOrdinarySQLShortScanPreservesReceiver(t *testing.T) {
	period, err := instant.Range(time.Unix(0, 0).UTC(), time.Unix(3600, 0).UTC())
	if err != nil {
		t.Fatal(err)
	}
	instantRange, err := postgres.NewInstantRange(period)
	if err != nil {
		t.Fatal(err)
	}
	dates, err := dateperiod.New(calendar.MustDate(2026, time.January, 1), calendar.MustDate(2026, time.January, 2), temporal.ClosedOpen)
	if err != nil {
		t.Fatal(err)
	}
	dateRange, err := postgres.NewDateRange(dates)
	if err != nil {
		t.Fatal(err)
	}
	for name, scanner := range map[string]interface {
		Scan(any) error
		Value() (driver.Value, error)
	}{"instant": &instantRange, "date": &dateRange} {
		t.Run(name, func(t *testing.T) {
			before, valueErr := scanner.Value()
			if valueErr != nil {
				t.Fatal(valueErr)
			}
			for _, input := range []string{"", "()"} {
				if scanErr := scanner.Scan(input); !errors.Is(scanErr, temporal.ErrUnsupported) {
					t.Fatalf("short scan = %v; want ErrUnsupported", scanErr)
				}
				after, valueErr := scanner.Value()
				if valueErr != nil || after != before {
					t.Fatalf("failed scan changed receiver: %v, %v; want %v", after, valueErr, before)
				}
			}
		})
	}
}

func TestOrdinaryCollectionFormatAndParseAdmission(t *testing.T) {
	period, err := instant.Range(time.Unix(0, 0).UTC(), time.Unix(3600, 0).UTC())
	if err != nil {
		t.Fatal(err)
	}
	instants, err := instant.NewSet(temporal.Limits{}, period)
	if err != nil {
		t.Fatal(err)
	}
	periodDates, err := dateperiod.New(calendar.MustDate(2026, time.January, 1), calendar.MustDate(2026, time.January, 2), temporal.ClosedOpen)
	if err != nil {
		t.Fatal(err)
	}
	dates, err := dateperiod.NewSet(temporal.Limits{}, periodDates)
	if err != nil {
		t.Fatal(err)
	}
	interval, err := timeofday.Between(timeofday.Midnight(), timeofday.Noon(), temporal.ClosedOpen)
	if err != nil {
		t.Fatal(err)
	}
	daily, err := timeofday.NewIntervalSet(temporal.Limits{}, interval)
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name         string
		kind         wire.Kind
		makeDocument func(temporal.Limits) (wire.CollectionDocument, error)
	}{
		{"instant", wire.KindInstantSet, func(l temporal.Limits) (wire.CollectionDocument, error) { return wire.FromInstantSet(instants, l) }},
		{"date", wire.KindDateSet, func(l temporal.Limits) (wire.CollectionDocument, error) { return wire.FromDateSet(dates, l) }},
		{"daily", wire.KindDailySet, func(l temporal.Limits) (wire.CollectionDocument, error) { return wire.FromDailySet(daily, l) }},
	} {
		t.Run(test.name, func(t *testing.T) {
			for _, limits := range []temporal.Limits{{FormatBytes: 1}, {ParseBytes: 1}} {
				got, admissionErr := test.makeDocument(limits)
				if !errors.Is(admissionErr, temporal.ErrLimit) || got.Version != "" || got.Kind != "" || got.Values != nil {
					t.Fatalf("limited document = %+v, %v; want zero and ErrLimit", got, admissionErr)
				}
			}
			got, controlErr := test.makeDocument(temporal.Limits{})
			if controlErr != nil || got.Version != wire.Version1 || got.Kind != test.kind || len(got.Values) != 1 {
				t.Fatalf("ordinary document = %+v, %v", got, controlErr)
			}
		})
	}
	if instants.Len() != 1 || dates.Len() != 1 || daily.Len() != 1 {
		t.Fatal("document admission changed input sets")
	}
}

func TestOrdinaryWireMalformedMembersPreservePrivateCategory(t *testing.T) {
	for _, input := range []string{
		`{"private":"application-private"]`,
		`[true}`,
		`{"private":"application-private",}`,
		`{"private":}`,
		`[true,]`,
		`[{"private":}]`,
		`{"private":"application-private"`,
	} {
		for name, decode := range map[string]func([]byte) (bool, error){
			"scalar": func(data []byte) (bool, error) {
				got, err := wire.Unmarshal(data, temporal.Limits{})
				return got == (wire.Document{}), err
			},
			"collection": func(data []byte) (bool, error) {
				got, err := wire.UnmarshalCollection(data, temporal.Limits{})
				return got.Version == "" && got.Kind == "" && got.Values == nil, err
			},
		} {
			t.Run(name+"/"+input, func(t *testing.T) {
				zero, err := decode([]byte(input))
				if !zero || !errors.Is(err, temporal.ErrParse) {
					t.Fatalf("malformed document = zero:%v, %v; want zero and ErrParse", zero, err)
				}
				want := "temporal: parse error: " + name + " document syntax"
				if err.Error() != want || strings.Contains(err.Error(), "application-private") {
					t.Fatalf("diagnostic = %q; want safe category %q", err.Error(), want)
				}
			})
		}
	}
}

func TestOrdinarySetUnionAdmissionAndDailyUpperEndpoint(t *testing.T) {
	period, err := instant.Range(time.Unix(0, 0).UTC(), time.Unix(3600, 0).UTC())
	if err != nil {
		t.Fatal(err)
	}
	datePeriod, err := dateperiod.New(calendar.MustDate(2026, time.January, 1), calendar.MustDate(2026, time.January, 2), temporal.ClosedOpen)
	if err != nil {
		t.Fatal(err)
	}
	for _, cap := range []int{1, 2} {
		left, buildErr := instant.NewSet(temporal.Limits{InputPeriods: cap}, period)
		if buildErr != nil {
			t.Fatal(buildErr)
		}
		right, buildErr := instant.NewSet(temporal.Limits{}, period)
		if buildErr != nil {
			t.Fatal(buildErr)
		}
		got, unionErr := left.Union(right)
		assertOrdinaryUnionAdmission(t, cap, got.Len(), unionErr)
		if !left.Equal(right) || left.Len() != 1 || right.Len() != 1 {
			t.Fatal("instant union changed input sets")
		}
		dateLeft, buildErr := dateperiod.NewSet(temporal.Limits{InputPeriods: cap}, datePeriod)
		if buildErr != nil {
			t.Fatal(buildErr)
		}
		dateRight, buildErr := dateperiod.NewSet(temporal.Limits{}, datePeriod)
		if buildErr != nil {
			t.Fatal(buildErr)
		}
		dateGot, unionErr := dateLeft.Union(dateRight)
		assertOrdinaryUnionAdmission(t, cap, dateGot.Len(), unionErr)
		if !dateLeft.Equal(dateRight) || dateLeft.Len() != 1 || dateRight.Len() != 1 {
			t.Fatal("date union changed input sets")
		}
	}
	closed, err := timeofday.Between(timeofday.Midnight(), timeofday.Noon(), temporal.Closed)
	if err != nil {
		t.Fatal(err)
	}
	open, err := timeofday.Between(timeofday.Midnight(), timeofday.Noon(), temporal.ClosedOpen)
	if err != nil {
		t.Fatal(err)
	}
	got, err := timeofday.NewIntervalSet(temporal.Limits{InputPeriods: 2}, closed, open)
	if err != nil || got.Len() != 1 || !got.Includes(timeofday.Noon()) || !got.Intervals()[0].Equal(closed) {
		t.Fatalf("daily endpoint normalization = %+v, %v; want closed interval", got, err)
	}
}

func assertOrdinaryUnionAdmission(t *testing.T, cap, length int, err error) {
	t.Helper()
	if cap == 2 {
		if err != nil || length != 1 {
			t.Fatalf("inclusive union = len:%d, %v", length, err)
		}
		return
	}
	var limit *temporal.LimitError
	if !errors.Is(err, temporal.ErrLimit) || !errors.As(err, &limit) || length != 0 {
		t.Fatalf("limited union = len:%d, %v; want zero and LimitError", length, err)
	}
	if limit.Field != "input_periods" || limit.Value != 2 || limit.Max != 1 {
		t.Fatalf("limit details = %+v; want input_periods=2 maximum1", limit)
	}
}
