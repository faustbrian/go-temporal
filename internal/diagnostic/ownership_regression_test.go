package diagnostic

import (
	"errors"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestDiagnosticTextDoesNotRetainInputBacking(t *testing.T) {
	message := "a€" + strings.Repeat("private-discarded-tail", 1024)
	sentinel := errors.New("safe classification")
	for _, test := range []struct {
		message string
		budget  int
		want    string
	}{
		{message, 1, "a"},
		{message, 3, "a"},
		{message, 4, "a€"},
		{message[:4], 4, "a€"},
	} {
		err := New(test.budget, test.message, sentinel)
		text := err.Error()
		if text != test.want || !utf8.ValidString(text) || !errors.Is(err, sentinel) {
			t.Fatalf("bounded diagnostic lost text/classification: %q, %v", text, err)
		}
		// Go's supported reflect string pointer observes actual backing ownership,
		// not allocation counts or a GC/finalizer timing proxy.
		if reflect.ValueOf(text).Pointer() == reflect.ValueOf(test.message).Pointer() {
			t.Error("bounded diagnostic retains the input string backing")
		}
		runtime.KeepAlive(test.message)
	}
}
