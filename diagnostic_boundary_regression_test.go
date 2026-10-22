package temporal_test

import (
	"errors"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	temporal "github.com/faustbrian/go-temporal/v3"
	config "github.com/faustbrian/go-temporal/v3/adapters/config"
	wire "github.com/faustbrian/go-temporal/v3/adapters/wire"
	"github.com/faustbrian/go-temporal/v3/instant"
	"github.com/faustbrian/go-temporal/v3/notation"
)

func TestRejectedEndpointDiagnosticsPreserveSafeClassificationAndOwnership(t *testing.T) {
	input := "private-temporal-input-" + strings.Repeat("z", 1024)
	for _, budget := range []int{1, 7, 1024} {
		period, err := notation.ParseInstant(input+"/2026-10-14T12:00:00Z", notation.ISO8601, temporal.Limits{ErrorBytes: budget})
		if err == nil || !period.IsEmpty() || !errors.Is(err, temporal.ErrParse) || len(err.Error()) > budget || !utf8.ValidString(err.Error()) {
			t.Fatalf("diagnostic budget %d: period=%+v err=%v", budget, period, err)
		}
		var unsafeCause *time.ParseError
		if errors.As(err, &unsafeCause) {
			t.Fatal("raw time parser cause exposed")
		}
		graph, ok := err.(interface{ Unwrap() []error })
		if !ok {
			t.Fatalf("missing stable cause graph: %T", err)
		}
		causes := graph.Unwrap()
		if len(causes) != 1 || !errors.Is(causes[0], temporal.ErrParse) || causes[0].Error() != temporal.ErrParse.Error() {
			t.Fatalf("unexpected unsafe causes: %v", causes)
		}
		for _, cause := range causes {
			if strings.Contains(cause.Error(), "private-temporal-input") {
				t.Fatal("rejected input retained in inspectable cause")
			}
		}
		causes[0] = errors.New("caller changed returned slice")
		if !errors.Is(err, temporal.ErrParse) || graph.Unwrap()[0].Error() != temporal.ErrParse.Error() {
			t.Fatal("caller changed published classification through cause slice")
		}
	}
	_, err := notation.ParseInstant(input, notation.ISO8601, temporal.Limits{ParseBytes: 10, ErrorBytes: 1})
	var limit *temporal.LimitError
	if !errors.Is(err, temporal.ErrLimit) || !errors.As(err, &limit) || limit.Field != "parse_bytes" || limit.Value != len(input) || limit.Max != 10 {
		t.Fatalf("safe structured admission cause lost: %v", err)
	}
}

func TestEndpointAdmissionThroughConfigurationAndWireDoesNotPublishPartialValue(t *testing.T) {
	original := instant.Point(time.Date(2026, time.October, 13, 12, 0, 0, 0, time.UTC))
	for _, endpoint := range []string{"2026-10-13T1:00:00Z", "2026-10-13T12:00:00+24:00", "2026-10-13T12:00:00.1234567890Z"} {
		value := config.NewInstantPeriod(original)
		text := "[" + endpoint + ",2026-10-15T12:00:00Z)"
		if err := value.UnmarshalText([]byte(text)); !errors.Is(err, temporal.ErrParse) || !value.Value().SetEqual(original) {
			t.Errorf("config endpoint admission published state: value=%+v err=%v", value.Value(), err)
		}
		document := wire.Document{Version: wire.Version1, Kind: wire.KindInstantPeriod, Value: text}
		period, err := document.Instant(temporal.Limits{})
		if !errors.Is(err, temporal.ErrParse) || !period.IsEmpty() {
			t.Errorf("wire endpoint admission = %+v, %v", period, err)
		}
	}
}
