package diagnostic

import (
	"errors"
	"testing"
	"unicode/utf8"
)

func TestErrorBoundsUTF8AndProtectsCauseSlice(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("sentinel")
	err := New(4, "a€b", sentinel)
	if err.Error() != "a€" || len(err.Error()) != 4 || !utf8.ValidString(err.Error()) ||
		!errors.Is(err, sentinel) {
		t.Fatalf("bounded error = %q", err)
	}

	var bounded *Error
	if !errors.As(err, &bounded) {
		t.Fatalf("bounded error type = %T", err)
	}
	causes := bounded.Unwrap()
	causes[0] = errors.New("changed")
	if !errors.Is(bounded, sentinel) {
		t.Fatal("caller mutated the published cause graph")
	}
}

func TestErrorHandlesEmptyAndAlreadyBoundedText(t *testing.T) {
	t.Parallel()

	if got := New(0, "message").Error(); got != "" {
		t.Fatalf("zero budget = %q", got)
	}
	if got := New(7, "message").Error(); got != "message" {
		t.Fatalf("bounded message = %q", got)
	}
	if got := New(3, "a€b").Error(); got != "a" {
		t.Fatalf("UTF-8 cut = %q", got)
	}
}
