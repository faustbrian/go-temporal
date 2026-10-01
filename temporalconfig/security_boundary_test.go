package temporalconfig_test

import (
	"errors"
	"testing"

	temporal "github.com/faustbrian/go-temporal/v2"
	"github.com/faustbrian/go-temporal/v2/temporalconfig"
)

func TestRetainedTextAdmissionRejectsNilReceiver(t *testing.T) {
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Error("nil receiver panicked instead of returning ErrUnsupported")
		}
	}()
	if err := (*temporalconfig.Time)(nil).UnmarshalText([]byte("08:00")); !errors.Is(err, temporal.ErrUnsupported) {
		t.Fatalf("nil receiver error = %v", err)
	}
}
