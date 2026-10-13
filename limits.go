package temporal

const (
	// HardMaxPeriods bounds raw input periods and normalized output periods.
	// Daily set Union permits twice InputPeriods normalized input segments
	// because each raw circular interval can expand into two segments.
	HardMaxPeriods = 100_000
	// HardMaxSteps bounds iteration and splitting.
	HardMaxSteps = 1_000_000

	hardMaxParseBytes  = 64 * 1024
	hardMaxPrecision   = 9
	hardMaxErrorBytes  = 1024
	hardMaxFormatBytes = 64 * 1024
	hardMaxParserDepth = 8
)

// Limits controls allocation and work performed by variable-size operations.
// A zero field selects its default; negative or over-hard-cap values fail
// Validate. Applicability is operation-specific: these are not a uniform total
// allocation, CPU, callback, precision-on-format or process-memory guarantee.
type Limits struct {
	// ParseBytes caps admitted input bytes at configured parser/text boundaries;
	// zero selects 64 KiB.
	ParseBytes int
	// Precision caps counted fractional-second digits in parsers; zero selects nine,
	// not zero digits. Formatters do not enforce this parse setting.
	Precision int
	// ErrorBytes caps selected diagnostic display text; zero selects 1024. It does
	// not cap every directly constructed error/cause graph.
	ErrorBytes int
	// FormatBytes caps selected emitted notation/JSON byte lengths; zero selects 64
	// KiB. Some formatters check after building text.
	FormatBytes int
	// InputPeriods caps raw constructor/collection input counts; zero selects
	// 100000. Daily Union uses twice this count for normalized segments.
	InputPeriods int
	// OutputPeriods caps normalized set outputs and instant splitting; not every
	// Steps/Split helper applies it. Zero selects 100000.
	OutputPeriods int
	// Steps caps selected iteration/split result counts; zero selects 1000000. It is
	// not a global callback/time budget.
	Steps int
	// ParserDepth caps nested JSON object/array depth at wire boundaries, with the
	// root at one; zero selects eight.
	ParserDepth int
}

// DefaultLimits returns the package's bounded operational defaults.
func DefaultLimits() Limits {
	return Limits{
		ParseBytes:    hardMaxParseBytes,
		Precision:     hardMaxPrecision,
		ErrorBytes:    hardMaxErrorBytes,
		FormatBytes:   hardMaxFormatBytes,
		InputPeriods:  HardMaxPeriods,
		OutputPeriods: HardMaxPeriods,
		Steps:         HardMaxSteps,
		ParserDepth:   hardMaxParserDepth,
	}
}

// Resolve replaces zero fields with defaults without mutation or validation.
// Negative and over-cap fields remain for Validate to reject.
func (l Limits) Resolve() Limits {
	defaults := DefaultLimits()

	if l.ParseBytes == 0 {
		l.ParseBytes = defaults.ParseBytes
	}
	if l.Precision == 0 {
		l.Precision = defaults.Precision
	}
	if l.ErrorBytes == 0 {
		l.ErrorBytes = defaults.ErrorBytes
	}
	if l.FormatBytes == 0 {
		l.FormatBytes = defaults.FormatBytes
	}
	if l.InputPeriods == 0 {
		l.InputPeriods = defaults.InputPeriods
	}
	if l.OutputPeriods == 0 {
		l.OutputPeriods = defaults.OutputPeriods
	}
	if l.Steps == 0 {
		l.Steps = defaults.Steps
	}
	if l.ParserDepth == 0 {
		l.ParserDepth = defaults.ParserDepth
	}

	return l
}

// Validate rejects negative fields and values above package hard limits; zero
// placeholders are accepted. It returns the first field violation as LimitError
// and does not resolve fields or validate any particular operation's input.
func (l Limits) Validate() error {
	checks := [...]struct {
		name  string
		value int
		max   int
	}{
		{"parse_bytes", l.ParseBytes, hardMaxParseBytes},
		{"precision", l.Precision, hardMaxPrecision},
		{"error_bytes", l.ErrorBytes, hardMaxErrorBytes},
		{"format_bytes", l.FormatBytes, hardMaxFormatBytes},
		{"input_periods", l.InputPeriods, HardMaxPeriods},
		{"output_periods", l.OutputPeriods, HardMaxPeriods},
		{"steps", l.Steps, HardMaxSteps},
		{"parser_depth", l.ParserDepth, hardMaxParserDepth},
	}

	for _, check := range checks {
		if check.value < 0 || check.value > check.max {
			return &LimitError{Field: check.name, Value: check.value, Max: check.max}
		}
	}

	return nil
}
