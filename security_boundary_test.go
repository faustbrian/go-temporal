package temporal_test

import (
	"errors"
	"testing"
	"time"

	temporal "github.com/faustbrian/go-temporal/v3"
	config "github.com/faustbrian/go-temporal/v3/adapters/config"
	wire "github.com/faustbrian/go-temporal/v3/adapters/wire"
	"github.com/faustbrian/go-temporal/v3/instant"
)

func TestTextAdmissionRejectsNilReceivers(t *testing.T) {
	for name, decode := range map[string]func([]byte) error{
		"canonical time": (*config.Time)(nil).UnmarshalText,
		"bounds":         func([]byte) error { return (*temporal.Bounds)(nil).UnmarshalText([]byte("[)")) },
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recovered := recover(); recovered != nil {
					t.Error("nil receiver panicked instead of returning ErrUnsupported")
				}
			}()
			if err := decode([]byte("08:00")); !errors.Is(err, temporal.ErrUnsupported) {
				t.Fatalf("nil receiver error = %v", err)
			}
		})
	}
}

func TestCollectionAdmissionPreservesInclusiveCountAndOwnedSet(t *testing.T) {
	base := time.Unix(0, 0).UTC()
	first, err := instant.Range(base, base.Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	second, err := instant.Range(base.Add(2*time.Hour), base.Add(3*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	set, err := instant.NewSet(temporal.Limits{}, first, second)
	if err != nil {
		t.Fatal(err)
	}
	accepted, err := wire.FromInstantSet(set, temporal.Limits{InputPeriods: 2})
	if err != nil || len(accepted.Values) != 2 {
		t.Fatalf("inclusive count = %+v, %v", accepted, err)
	}
	rejected, err := wire.FromInstantSet(set, temporal.Limits{InputPeriods: 1})
	if !errors.Is(err, temporal.ErrLimit) || rejected.Version != "" || rejected.Kind != "" || len(rejected.Values) != 0 {
		t.Fatalf("over count = %+v, %v", rejected, err)
	}
	if set.Len() != 2 || !set.Periods()[0].SetEqual(first) || !set.Periods()[1].SetEqual(second) {
		t.Fatal("collection admission changed the source set")
	}
}

func TestCollectionMarshalPreservesNilAndEmptyBudgetBoundaries(t *testing.T) {
	for name, encode := range map[string]func(wire.CollectionDocument, temporal.Limits) ([]byte, error){
		"canonical": wire.MarshalCollection,
	} {
		t.Run(name, func(t *testing.T) {
			for _, test := range []struct {
				name   string
				values []string
				want   string
			}{
				{"nil", nil, `{"version":"temporal/v1","kind":"instant-set","values":null}`},
				{"empty", []string{}, `{"version":"temporal/v1","kind":"instant-set","values":null}`},
			} {
				t.Run(test.name, func(t *testing.T) {
					document := wire.CollectionDocument{Version: wire.Version1, Kind: wire.KindInstantSet, Values: test.values}
					got, err := encode(document, temporal.Limits{FormatBytes: len(test.want)})
					if err != nil || string(got) != test.want {
						t.Fatalf("inclusive budget = %s, %v; want %s", got, err, test.want)
					}
					got, err = encode(document, temporal.Limits{FormatBytes: len(test.want) - 1})
					if !errors.Is(err, temporal.ErrLimit) || got != nil {
						t.Fatalf("over budget = %s, %v; want nil and ErrLimit", got, err)
					}
					got, err = encode(document, temporal.Limits{FormatBytes: len(`{"version":"temporal/v1","kind":"instant-set","values":[]}`)})
					if !errors.Is(err, temporal.ErrLimit) || got != nil {
						t.Fatalf("array-sized budget = %s, %v; want normalized null rejection", got, err)
					}
					if (document.Values == nil) != (test.values == nil) || len(document.Values) != 0 {
						t.Fatal("marshaling changed the caller-owned empty collection")
					}
				})
			}
		})
	}
}
