// Package diagnostic bounds display text and snapshots the outer cause slice.
// New does not sanitize, deep-copy, authenticate or redact arbitrary causes; its
// internal callers must select safe static/typed causes before publication.
// It owns no IO and provides no erasure or independent retained-memory budget.
package diagnostic

import "unicode/utf8"

// Error holds immutable bounded display text and an outer cause snapshot.
// Cause objects remain shared and caller-owned; their mutation and text are not
// bounded by this type. Nil receiver methods are unsupported.
type Error struct {
	message string
	causes  []error
}

// New truncates message to at most maxBytes, returning empty display text for
// non-positive limits. A truncated prefix is shortened to valid UTF-8; admitted
// shorter strings are not revalidated. Prefix/string backing and cause objects
// may remain shared, so callers must supply safe text/causes. No input is mutated.
func New(maxBytes int, message string, causes ...error) error {
	return &Error{
		message: truncateUTF8(message, maxBytes),
		causes:  append([]error(nil), causes...),
	}
}

// Error returns the stored bounded display text; it does not format causes.
func (e *Error) Error() string { return e.message }

// Unwrap returns a new outer slice. Cause objects remain borrowed; this does
// not prevent mutation through an aliased cause or bound its nested graph.
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
