package temporalwire_test

import (
	"errors"
	"strings"
	"testing"

	temporal "github.com/faustbrian/go-temporal/v3"
	wire "github.com/faustbrian/go-temporal/v3/adapters/wire"
)

func TestWireAdmissionRejectsDuplicateFieldsAndConfiguredDepth(t *testing.T) {
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
			got, err := wire.Unmarshal([]byte(test.input), test.limits)
			if !errors.Is(err, test.want) || got != (wire.Document{}) {
				t.Fatalf("Unmarshal = %+v, %v; want zero document and %v", got, err, test.want)
			}
			if strings.Contains(err.Error(), "application-private") {
				t.Fatal("diagnostic contains private input")
			}
		})
	}
	got, err := wire.Unmarshal([]byte(`{"version":"temporal/v1","kind":"time-of-day","value":"08:00"}`), temporal.Limits{ParserDepth: 1})
	if err != nil || got.Value != "08:00" {
		t.Fatalf("inclusive depth = %+v, %v", got, err)
	}
}

func TestCollectionWireAdmissionRejectsDuplicateFieldsAndConfiguredDepth(t *testing.T) {
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
			got, err := wire.UnmarshalCollection([]byte(test.input), test.limits)
			if !errors.Is(err, test.want) || got.Version != "" || got.Kind != "" || len(got.Values) != 0 {
				t.Fatalf("UnmarshalCollection = %+v, %v; want zero collection and %v", got, err, test.want)
			}
		})
	}
	got, err := wire.UnmarshalCollection([]byte(`{"version":"temporal/v1","kind":"instant-set","values":[]}`), temporal.Limits{ParserDepth: 2})
	if err != nil || got.Kind != wire.KindInstantSet || len(got.Values) != 0 {
		t.Fatalf("inclusive depth = %+v, %v", got, err)
	}
}
