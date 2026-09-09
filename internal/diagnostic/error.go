package diagnostic

import "unicode/utf8"

// Error is an immutable bounded diagnostic with a safe cause set.
type Error struct {
	message string
	causes  []error
}

// New returns a diagnostic whose text is bounded without retaining source text.
func New(maxBytes int, message string, causes ...error) error {
	return &Error{
		message: truncateUTF8(message, maxBytes),
		causes:  append([]error(nil), causes...),
	}
}

func (e *Error) Error() string { return e.message }

// Unwrap returns a copy so callers cannot mutate the published cause graph.
func (e *Error) Unwrap() []error { return append([]error(nil), e.causes...) }

func truncateUTF8(value string, maxBytes int) string {
	if maxBytes <= 0 {
		return ""
	}
	if len(value) <= maxBytes {
		return value
	}
	value = value[:maxBytes]
	for !utf8.ValidString(value) {
		value = value[:len(value)-1]
	}
	return value
}
