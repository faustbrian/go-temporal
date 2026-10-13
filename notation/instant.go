// Package notation implements pinned interval/fixed-duration profiles, not full
// ISO standards. ParseInstant delegates endpoints to Go time.RFC3339Nano after
// byte/UTF-8 and dot-fraction checks; Date and local-time parsing have their own
// exact grammars. ISO8601 start/end implies ClosedOpen; ISO80000 and Bourbaki
// are explicit bracket alphabets. See docs/specification-decisions.md (RFC3339
// and pinned compatibility profile 469603239dbe700739c29b4c532a90382b6cbedf).
//
// Limits zero fields resolve to defaults. Formatting preserves selected bounds
// but is not general revalidation or a same-limit parse round-trip promise:
// precision is a parse setting, output checks can follow text construction, and
// in-memory values can lie outside endpoint parser domains. Calls borrow inputs,
// return new text/values or zero results on failure, and retain no IO resources.
// Pinned Go parser permissiveness is not certification of the desired strict
// profile; normative decisions and conformance expectations remain separate.
package notation

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	calendar "github.com/faustbrian/go-calendar/v2"
	temporal "github.com/faustbrian/go-temporal/v2"
	"github.com/faustbrian/go-temporal/v2/instant"
	"github.com/faustbrian/go-temporal/v2/internal/diagnostic"
)

// Format identifies a supported interval notation.
type Format uint8

const (
	// FormatInvalid is the safe zero value.
	FormatInvalid Format = iota
	// ISO8601 uses start/end form and implies closed-open bounds.
	ISO8601
	// ISO80000 uses conventional inclusive and exclusive brackets.
	ISO80000
	// Bourbaki uses inward and outward square brackets.
	Bourbaki
)

// ParseInstant decodes one complete interval with resolved byte/UTF-8/bounds
// checks and zero Period on failure. Dot-fraction digit counting precedes the
// pinned Go RFC3339Nano parser; other forms accepted by that parser are not all
// covered by this local precision check. This existing admission limitation is
// not a strict-RFC conformance guarantee or an intended policy interpretation.
func ParseInstant(value string, format Format, limits temporal.Limits) (instant.Period, error) {
	limits = limits.Resolve()
	if err := limits.Validate(); err != nil {
		return instant.Period{}, err
	}
	if len(value) > limits.ParseBytes {
		cause := &temporal.LimitError{
			Field: "parse_bytes",
			Value: len(value),
			Max:   limits.ParseBytes,
		}
		return instant.Period{}, diagnostic.New(limits.ErrorBytes, temporal.ErrLimit.Error(), temporal.ErrLimit, cause)
	}
	if !utf8.ValidString(value) {
		return instant.Period{}, diagnostic.New(limits.ErrorBytes, "temporal: parse error: syntax", temporal.ErrParse)
	}

	var startText, endText string
	var bounds temporal.Bounds
	var err error
	switch format {
	case ISO8601:
		startText, endText, err = splitExactly(value, '/')
		bounds = temporal.ClosedOpen
	case ISO80000:
		startText, endText, bounds, err = splitBounded(value, false)
	case Bourbaki:
		startText, endText, bounds, err = splitBounded(value, true)
	default:
		return instant.Period{}, diagnostic.New(limits.ErrorBytes, temporal.ErrUnsupported.Error(), temporal.ErrUnsupported)
	}
	if err != nil {
		return instant.Period{}, boundedParseError(limits, "interval syntax", err)
	}

	start, err := parseInstant(startText, limits.Precision)
	if err != nil {
		return instant.Period{}, boundedParseError(limits, "start timestamp", err)
	}
	end, err := parseInstant(endText, limits.Precision)
	if err != nil {
		return instant.Period{}, boundedParseError(limits, "end timestamp", err)
	}

	period, err := instant.New(start, end, bounds)
	if err != nil {
		return instant.Period{}, boundedParseError(limits, "interval bounds", err)
	}

	return period, nil
}

func boundedParseError(limits temporal.Limits, stage string, cause error) error {
	causes := []error{temporal.ErrParse}
	for _, sentinel := range []error{temporal.ErrBounds, temporal.ErrPrecision, temporal.ErrInvalidTime,
		temporal.ErrUnsupported, temporal.ErrOverflow, temporal.ErrReversed} {
		if errors.Is(cause, sentinel) {
			causes = append(causes, sentinel)
		}
	}
	for _, sentinel := range []error{calendar.ErrInvalidFormat, calendar.ErrInvalidDate} {
		if errors.Is(cause, sentinel) {
			causes = append(causes, sentinel)
		}
	}
	return diagnostic.New(limits.ErrorBytes, "temporal: parse error: "+stage, causes...)
}

// FormatInstant emits declared endpoints in Go RFC3339Nano spelling and the
// selected bounds profile; ISO8601 rejects non-ClosedOpen bounds. Location
// identity is not encoded. It does not reparse endpoints or enforce parse
// Precision/year-domain validity; FormatBytes is checked after text assembly.
// Returns empty text on error, not a universal same-limit round-trip promise.
func FormatInstant(period instant.Period, format Format, limits temporal.Limits) (string, error) {
	limits = limits.Resolve()
	if err := limits.Validate(); err != nil {
		return "", err
	}

	start := period.Start().Format(time.RFC3339Nano)
	end := period.End().Format(time.RFC3339Nano)
	var value string

	switch format {
	case ISO8601:
		if period.Bounds() != temporal.ClosedOpen {
			return "", fmt.Errorf("%w: ISO 8601 start/end does not encode bounds", temporal.ErrUnsupported)
		}
		value = start + "/" + end
	case ISO80000:
		left, right := isoBrackets(period.Bounds())
		value = left + start + "," + end + right
	case Bourbaki:
		left, right := bourbakiBrackets(period.Bounds())
		value = left + start + "," + end + right
	default:
		return "", temporal.ErrUnsupported
	}

	if len(value) > limits.FormatBytes {
		return "", &temporal.LimitError{
			Field: "format_bytes",
			Value: len(value),
			Max:   limits.FormatBytes,
		}
	}

	return value, nil
}

func splitBounded(value string, bourbaki bool) (string, string, temporal.Bounds, error) {
	if len(value) < 5 {
		return "", "", temporal.ClosedOpen, temporal.ErrParse
	}

	left := value[0]
	right := value[len(value)-1]
	includeStart, includeEnd, ok := bracketInclusion(left, right, bourbaki)
	if !ok {
		return "", "", temporal.ClosedOpen, temporal.ErrBounds
	}
	start, end, err := splitExactly(value[1:len(value)-1], ',')
	if err != nil {
		return "", "", temporal.ClosedOpen, err
	}

	return start, end, notationBounds(includeStart, includeEnd), nil
}

func splitExactly(value string, separator byte) (string, string, error) {
	index := strings.IndexByte(value, separator)
	if index <= 0 || index == len(value)-1 || strings.IndexByte(value[index+1:], separator) >= 0 {
		return "", "", temporal.ErrParse
	}

	return value[:index], value[index+1:], nil
}

func bracketInclusion(left, right byte, bourbaki bool) (bool, bool, bool) {
	if bourbaki {
		if (left != '[' && left != ']') || (right != '[' && right != ']') {
			return false, false, false
		}
		return left == '[', right == ']', true
	}

	if (left != '[' && left != '(') || (right != ']' && right != ')') {
		return false, false, false
	}
	return left == '[', right == ']', true
}

func notationBounds(includeStart, includeEnd bool) temporal.Bounds {
	switch {
	case includeStart && includeEnd:
		return temporal.Closed
	case includeStart:
		return temporal.ClosedOpen
	case includeEnd:
		return temporal.OpenClosed
	default:
		return temporal.Open
	}
}

func isoBrackets(bounds temporal.Bounds) (string, string) {
	left := "("
	right := ")"
	if bounds.IncludesStart() {
		left = "["
	}
	if bounds.IncludesEnd() {
		right = "]"
	}

	return left, right
}

func bourbakiBrackets(bounds temporal.Bounds) (string, string) {
	left := "]"
	right := "["
	if bounds.IncludesStart() {
		left = "["
	}
	if bounds.IncludesEnd() {
		right = "]"
	}

	return left, right
}

func parseInstant(value string, precision int) (time.Time, error) {
	if fractionDigits(value) > precision {
		return time.Time{}, temporal.ErrPrecision
	}

	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, err
	}

	return parsed, nil
}

func fractionDigits(value string) int {
	dot := strings.IndexByte(value, '.')
	if dot < 0 {
		return 0
	}

	digits := 0
	for index := dot + 1; index < len(value) && value[index] >= '0' && value[index] <= '9'; index++ {
		digits++
	}

	return digits
}
