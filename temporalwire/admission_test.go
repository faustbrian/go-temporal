package temporalwire_test

import (
	"errors"
	"testing"

	temporal "github.com/faustbrian/go-temporal/v2"
	"github.com/faustbrian/go-temporal/v2/temporalwire"
)

func TestRetainedCollectionAdmissionPreservesCallerAndOwnedResults(t *testing.T) {
	// These are public outcome/ownership controls. The removal of a pre-admission
	// copy is established at the adapter owner, not by an allocation assertion.
	for _, test := range []struct {
		name   string
		kind   temporalwire.Kind
		values []string
		decode func(temporalwire.CollectionDocument, temporal.Limits) (func() (int, string, error), error)
	}{
		{"instant", temporalwire.KindInstantSet, []string{"[1970-01-01T00:00:00Z,1970-01-01T01:00:00Z)", "[1970-01-01T02:00:00Z,1970-01-01T03:00:00Z)"},
			func(document temporalwire.CollectionDocument, limits temporal.Limits) (func() (int, string, error), error) {
				set, err := document.InstantSet(limits)
				return func() (int, string, error) {
					snapshot, snapshotErr := temporalwire.FromInstantSet(set, temporal.Limits{})
					return set.Len(), firstCollectionValue(snapshot), snapshotErr
				}, err
			}},
		{"date", temporalwire.KindDateSet, []string{"[2026-01-01,2026-01-01]", "[2026-01-03,2026-01-03]"},
			func(document temporalwire.CollectionDocument, limits temporal.Limits) (func() (int, string, error), error) {
				set, err := document.DateSet(limits)
				return func() (int, string, error) {
					snapshot, snapshotErr := temporalwire.FromDateSet(set, temporal.Limits{})
					return set.Len(), firstCollectionValue(snapshot), snapshotErr
				}, err
			}},
		{"daily", temporalwire.KindDailySet, []string{"[08:00:00.000000000,09:00:00.000000000)", "[10:00:00.000000000,11:00:00.000000000)"},
			func(document temporalwire.CollectionDocument, limits temporal.Limits) (func() (int, string, error), error) {
				set, err := document.DailySet(limits)
				return func() (int, string, error) {
					snapshot, snapshotErr := temporalwire.FromDailySet(set, temporal.Limits{})
					return set.Len(), firstCollectionValue(snapshot), snapshotErr
				}, err
			}},
	} {
		t.Run(test.name, func(t *testing.T) {
			original := append([]string(nil), test.values...)
			document := temporalwire.CollectionDocument{Version: temporalwire.Version1, Kind: test.kind, Values: test.values}
			failedSnapshot, err := test.decode(document, temporal.Limits{InputPeriods: 1})
			if !errors.Is(err, temporal.ErrLimit) {
				t.Fatalf("over-count decode error = %v", err)
			}
			count, first, err := failedSnapshot()
			if err != nil || count != 0 || first != "" {
				t.Fatalf("failed set snapshot = %d, %q, %v", count, first, err)
			}
			encoded, err := temporalwire.MarshalCollection(document, temporal.Limits{InputPeriods: 1})
			if !errors.Is(err, temporal.ErrLimit) || encoded != nil {
				t.Fatalf("over-count marshal = %s, %v", encoded, err)
			}
			ownedSnapshot, err := test.decode(document, temporal.Limits{InputPeriods: 2})
			if err != nil {
				t.Fatal(err)
			}
			encoded, err = temporalwire.MarshalCollection(document, temporal.Limits{InputPeriods: 2})
			if err != nil || len(encoded) == 0 || document.Values[0] != original[0] || document.Values[1] != original[1] {
				t.Fatalf("inclusive marshal or caller ownership changed: %v", err)
			}
			decoded, err := temporalwire.UnmarshalCollection(encoded, temporal.Limits{InputPeriods: 2})
			if err != nil || len(decoded.Values) != 2 {
				t.Fatalf("marshal snapshot = %+v, %v", decoded, err)
			}
			document.Values[0] = "caller replacement"
			count, first, err = ownedSnapshot()
			if err != nil || count != 2 || first != original[0] || decoded.Values[0] != original[0] {
				t.Fatalf("owned snapshot aliased caller: %d, %q, %v", count, first, err)
			}
		})
	}
}

func firstCollectionValue(document temporalwire.CollectionDocument) string {
	if len(document.Values) == 0 {
		return ""
	}
	return document.Values[0]
}
