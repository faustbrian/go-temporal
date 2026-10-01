package temporalwire_test

import (
	"errors"
	"strings"
	"testing"

	temporal "github.com/faustbrian/go-temporal/v2"
	wire "github.com/faustbrian/go-temporal/v2/adapters/wire"
	legacywire "github.com/faustbrian/go-temporal/v2/temporalwire"
)

func TestRetainedWireAdmissionRejectsDuplicateFieldsAndConfiguredDepth(t *testing.T) {
	for name, decode := range map[string]func([]byte, temporal.Limits) (wire.Document, error){
		"retained": func(input []byte, limits temporal.Limits) (wire.Document, error) {
			value, err := legacywire.Unmarshal(input, limits)
			return wire.Document{Version: value.Version, Kind: wire.Kind(value.Kind), Value: value.Value}, err
		},
	} {
		t.Run(name, func(t *testing.T) {
			for _, test := range []struct {
				name   string
				input  string
				limits temporal.Limits
				want   error
			}{
				{"duplicate", `{"version":"temporal/v1","kind":"time-of-day","value":"08:00","value":"09:00"}`, temporal.Limits{}, temporal.ErrParse},
				{"depth", `{"version":"temporal/v1","kind":"time-of-day","value":"08:00","private":{"note":"application-private"}}`, temporal.Limits{ParserDepth: 1}, temporal.ErrLimit},
			} {
				t.Run(test.name, func(t *testing.T) {
					got, err := decode([]byte(test.input), test.limits)
					if !errors.Is(err, test.want) || got != (wire.Document{}) {
						t.Fatalf("decode = %+v, %v; want zero document and %v", got, err, test.want)
					}
					if strings.Contains(err.Error(), "application-private") {
						t.Fatal("diagnostic contains private input")
					}
				})
			}
			got, err := decode([]byte(`{"version":"temporal/v1","kind":"time-of-day","value":"08:00"}`), temporal.Limits{ParserDepth: 1})
			if err != nil || got.Value != "08:00" {
				t.Fatalf("inclusive depth = %+v, %v", got, err)
			}
		})
	}
}

func TestRetainedCollectionWireAdmissionRejectsDuplicateFieldsAndConfiguredDepth(t *testing.T) {
	for name, decode := range map[string]func([]byte, temporal.Limits) (wire.CollectionDocument, error){
		"retained": func(input []byte, limits temporal.Limits) (wire.CollectionDocument, error) {
			value, err := legacywire.UnmarshalCollection(input, limits)
			return wire.CollectionDocument{Version: value.Version, Kind: wire.Kind(value.Kind), Values: value.Values}, err
		},
	} {
		t.Run(name, func(t *testing.T) {
			for _, test := range []struct {
				name   string
				input  string
				limits temporal.Limits
				want   error
			}{
				{"duplicate", `{"version":"temporal/v1","kind":"instant-set","values":[],"values":[]}`, temporal.Limits{}, temporal.ErrParse},
				{"depth", `{"version":"temporal/v1","kind":"instant-set","values":[]}`, temporal.Limits{ParserDepth: 1}, temporal.ErrLimit},
			} {
				t.Run(test.name, func(t *testing.T) {
					got, err := decode([]byte(test.input), test.limits)
					if !errors.Is(err, test.want) || got.Version != "" || got.Kind != "" || len(got.Values) != 0 {
						t.Fatalf("decode = %+v, %v; want zero collection and %v", got, err, test.want)
					}
				})
			}
			got, err := decode([]byte(`{"version":"temporal/v1","kind":"instant-set","values":[]}`), temporal.Limits{ParserDepth: 2})
			if err != nil || got.Kind != wire.KindInstantSet || len(got.Values) != 0 {
				t.Fatalf("inclusive depth = %+v, %v", got, err)
			}
		})
	}
}

func TestRetainedCollectionMarshalPreservesNilAndEmptyBudgetBoundaries(t *testing.T) {
	for name, encode := range map[string]func(wire.CollectionDocument, temporal.Limits) ([]byte, error){
		"retained": func(document wire.CollectionDocument, limits temporal.Limits) ([]byte, error) {
			return legacywire.MarshalCollection(legacywire.CollectionDocument{
				Version: document.Version, Kind: legacywire.Kind(document.Kind), Values: document.Values,
			}, limits)
		},
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
