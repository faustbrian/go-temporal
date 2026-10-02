package temporalwire

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	temporal "github.com/faustbrian/go-temporal/v2"
)

func TestJSONStringSizeMatchesStandardEncodingWithoutChangingRepresentation(t *testing.T) {
	values := []string{"", " ", "08:00", "quoted\"\\", "\b\f\n\r\t\x01", "<>&", "é\u2028\u2029", "\ufffd"}
	completed := make(chan []int, 1)
	go func() {
		sizes := make([]int, len(values))
		for index, value := range values {
			sizes[index] = jsonStringSize(value)
		}
		completed <- sizes
	}()
	// The selected verifier grants each mutant at least one minute. A broken
	// termination guard must fail this ordinary oracle before that watchdog.
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()
	var sizes []int
	select {
	case sizes = <-completed:
	case <-timer.C:
		t.Fatal("JSON string size calculation did not complete within five seconds")
	}
	for index, value := range values {
		encoded, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		if got := sizes[index]; got != len(encoded) {
			t.Fatalf("encoded size = %d, want %d", got, len(encoded))
		}
	}
}

func TestCollectionCountAdmissionBeforeCopyingAndFormatting(t *testing.T) {
	if err := admitCollectionCount(2, temporal.Limits{InputPeriods: 2}); err != nil {
		t.Fatalf("inclusive count: %v", err)
	}
	for _, test := range []struct {
		name   string
		count  int
		limits temporal.Limits
		field  string
	}{
		{"count", 2, temporal.Limits{InputPeriods: 1}, "input_periods"},
		{"invalid limits", 0, temporal.Limits{ErrorBytes: -1}, "error_bytes"},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := admitCollectionCount(test.count, test.limits)
			var limitError *temporal.LimitError
			if !errors.Is(err, temporal.ErrLimit) || !errors.As(err, &limitError) || limitError.Field != test.field {
				t.Fatalf("admission = %v; want ErrLimit for %s", err, test.field)
			}
		})
	}
}

func TestCollectionSizeAdmissionPreservesNilAndEmptyEncoding(t *testing.T) {
	for name, values := range map[string][]string{"nil": nil, "empty": {}} {
		t.Run(name, func(t *testing.T) {
			document := CollectionDocument{Version: Version1, Kind: KindInstantSet, Values: values}
			encoded, err := json.Marshal(document)
			if err != nil {
				t.Fatal(err)
			}
			for _, maximum := range []int{len(encoded), len(encoded) - 1} {
				size, exceeds := collectionDocumentSize(document, maximum)
				if size != len(encoded) || exceeds != (maximum < len(encoded)) {
					t.Fatalf("size admission at %d = %d, %t; encoded size %d", maximum, size, exceeds, len(encoded))
				}
			}
		})
	}
}
