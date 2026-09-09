package notation

import (
	"fmt"
	"unicode/utf8"

	calendar "github.com/faustbrian/go-calendar"
	temporal "github.com/faustbrian/go-temporal"
	"github.com/faustbrian/go-temporal/dateperiod"
	"github.com/faustbrian/go-temporal/internal/diagnostic"
)

// ParseDate decodes one complete bounded civil-date interval.
func ParseDate(value string, format Format, limits temporal.Limits) (dateperiod.Period, error) {
	limits = limits.Resolve()
	if err := limits.Validate(); err != nil {
		return dateperiod.Period{}, err
	}
	if len(value) > limits.ParseBytes {
		cause := &temporal.LimitError{
			Field: "parse_bytes", Value: len(value), Max: limits.ParseBytes,
		}
		return dateperiod.Period{}, diagnostic.New(limits.ErrorBytes, temporal.ErrLimit.Error(), temporal.ErrLimit, cause)
	}
	if !utf8.ValidString(value) {
		return dateperiod.Period{}, diagnostic.New(limits.ErrorBytes, "temporal: parse error: syntax", temporal.ErrParse)
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
		return dateperiod.Period{}, diagnostic.New(limits.ErrorBytes, temporal.ErrUnsupported.Error(), temporal.ErrUnsupported)
	}
	if err != nil {
		return dateperiod.Period{}, boundedParseError(limits, "interval syntax", err)
	}

	start, err := calendar.ParseDate(startText)
	if err != nil {
		return dateperiod.Period{}, boundedParseError(limits, "start date", err)
	}
	end, err := calendar.ParseDate(endText)
	if err != nil {
		return dateperiod.Period{}, boundedParseError(limits, "end date", err)
	}
	period, err := dateperiod.New(start, end, bounds)
	if err != nil {
		return dateperiod.Period{}, boundedParseError(limits, "interval bounds", err)
	}
	return period, nil
}

// FormatDate encodes a bounded civil-date interval without semantic loss.
func FormatDate(period dateperiod.Period, format Format, limits temporal.Limits) (string, error) {
	limits = limits.Resolve()
	if err := limits.Validate(); err != nil {
		return "", err
	}

	start := period.Start().String()
	end := period.End().String()
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
			Field: "format_bytes", Value: len(value), Max: limits.FormatBytes,
		}
	}
	return value, nil
}
