package temporalwire_test

import (
	"errors"
	"slices"
	"strings"
	"testing"

	temporal "github.com/faustbrian/go-temporal/v2"
	wire "github.com/faustbrian/go-temporal/v2/temporalwire"
)

func TestWireScalarMarshalExactEncodedBudget(t *testing.T) {
	document := wire.Document{Version: wire.Version1, Kind: wire.KindTime, Value: "08:00"}
	const want = `{"version":"temporal/v1","kind":"time-of-day","value":"08:00"}`
	for name, marshal := range map[string]func(wire.Document, temporal.Limits) ([]byte, error){
		"retained": wire.Marshal,
	} {
		t.Run(name, func(t *testing.T) {
			got, err := marshal(document, temporal.Limits{FormatBytes: len(want)})
			if err != nil || string(got) != want {
				t.Fatalf("exact budget = %s, %v; want %s", got, err, want)
			}
			got, err = marshal(document, temporal.Limits{FormatBytes: len(want) - 1})
			if !errors.Is(err, temporal.ErrLimit) || got != nil {
				t.Fatalf("one-less budget = %s, %v; want nil and ErrLimit", got, err)
			}
			if document.Value != "08:00" || document.Version != wire.Version1 || document.Kind != wire.KindTime {
				t.Fatal("marshal changed caller document")
			}
		})
	}
}

func TestWireCollectionMarshalExactEncodedBudget(t *testing.T) {
	const first = "[01:00,02:00)"
	const second = "[03:00,04:00)"
	const want = `{"version":"temporal/v1","kind":"daily-set","values":["[01:00,02:00)","[03:00,04:00)"]}`
	for name, marshal := range map[string]func(wire.CollectionDocument, temporal.Limits) ([]byte, error){
		"retained": wire.MarshalCollection,
	} {
		t.Run(name, func(t *testing.T) {
			values := []string{first, second}
			document := wire.CollectionDocument{Version: wire.Version1, Kind: wire.KindDailySet, Values: values}
			got, err := marshal(document, temporal.Limits{FormatBytes: len(want), InputPeriods: 2})
			if err != nil || string(got) != want {
				t.Fatalf("exact budget = %s, %v; want %s", got, err, want)
			}
			got, err = marshal(document, temporal.Limits{FormatBytes: len(want) - 1, InputPeriods: 2})
			if !errors.Is(err, temporal.ErrLimit) || got != nil {
				t.Fatalf("one-less budget = %s, %v; want nil and ErrLimit", got, err)
			}
			if document.Version != wire.Version1 || document.Kind != wire.KindDailySet || !slices.Equal(values, []string{first, second}) {
				t.Fatal("marshal changed caller collection")
			}
		})
	}
}

func TestWireNestedArrayDepthPrecedesTypedContent(t *testing.T) {
	const input = `{"version":"temporal/v1","kind":"daily-set","values":[["application-private"]]}`
	for name, unmarshal := range map[string]func([]byte, temporal.Limits) (wire.CollectionDocument, error){
		"retained": wire.UnmarshalCollection,
	} {
		t.Run(name, func(t *testing.T) {
			got, err := unmarshal([]byte(input), temporal.Limits{ParserDepth: 2})
			if !errors.Is(err, temporal.ErrLimit) || got.Version != "" || got.Kind != "" || got.Values != nil {
				t.Fatalf("nested array = %+v, %v; want zero document and ErrLimit", got, err)
			}
			if strings.Contains(err.Error(), "application-private") {
				t.Fatal("depth diagnostic exposed private value")
			}
			got, err = unmarshal([]byte(`{"version":"temporal/v1","kind":"daily-set","values":[]}`), temporal.Limits{ParserDepth: 2})
			if err != nil || got.Version != wire.Version1 || got.Kind != wire.KindDailySet || len(got.Values) != 0 {
				t.Fatalf("inclusive depth = %+v, %v", got, err)
			}
		})
	}
}
