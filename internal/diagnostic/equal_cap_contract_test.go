package diagnostic_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/faustbrian/go-temporal/v2/internal/diagnostic"
)

type boundaryCause struct{ message string }

func (cause *boundaryCause) Error() string { return cause.message }

func TestDiagnosticConstructorPreservesAdmittedBytesAndCauseOwnership(t *testing.T) {
	malformed := string([]byte{0xff})
	for _, test := range []struct {
		name    string
		limit   int
		message string
		want    string
	}{
		{"equal valid ASCII", 1, "a", "a"},
		{"equal valid multibyte", 3, "€", "€"},
		{"equal malformed byte", 1, malformed, malformed},
		{"under malformed byte", 2, malformed, malformed},
		{"over valid boundary", 4, "a€b", "a€"},
		{"over incomplete rune", 3, "a€b", "a"},
		{"over malformed prefix", 1, malformed + "a", ""},
		{"over valid prefix", 1, "a" + malformed, "a"},
		{"zero nonempty", 0, malformed, ""},
		{"zero empty", 0, "", ""},
		{"negative", -1, malformed, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			cause := &boundaryCause{message: "synthetic classification"}
			inputs := []error{cause}
			result := diagnostic.New(test.limit, test.message, inputs...)
			if result.Error() != test.want {
				t.Errorf("diagnostic bytes = %x, want %x", []byte(result.Error()), []byte(test.want))
			}
			inputs[0] = errors.New("changed input slice")
			if !errors.Is(result, cause) {
				t.Error("constructor lost original cause identity after caller slice mutation")
			}
			var typed *boundaryCause
			if !errors.As(result, &typed) || typed != cause {
				t.Error("constructor lost the inspectable original typed cause")
			}
			var bounded *diagnostic.Error
			if !errors.As(result, &bounded) {
				t.Fatalf("constructor error type = %T", result)
			}
			published := bounded.Unwrap()
			// Unwrap snapshots retain the entries themselves, not just equivalent causes.
			if !slices.Equal(published, []error{cause}) {
				t.Fatal("unwrapped cause snapshot lost original identity")
			}
			published[0] = errors.New("changed returned slice")
			again := bounded.Unwrap()
			if !slices.Equal(again, []error{cause}) || !errors.Is(result, cause) {
				t.Error("caller mutation of Unwrap changed retained cause ownership")
			}
		})
	}
}
