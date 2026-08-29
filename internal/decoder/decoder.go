package decoder

import "strings"

type SerialFormat string

const (
	FormatLegacy11   SerialFormat = "Legacy 11-character"
	FormatLegacy12   SerialFormat = "Legacy 12-character"
	FormatRandomized SerialFormat = "Randomized"
	FormatInvalid    SerialFormat = "Invalid"
)

type Result struct {
	Serial       string
	Format       SerialFormat
	Valid        bool
	Location     string
	LocationInfo *Location
	YearCode     string
	Year         int
	Half         int
	WeekCode     string
	Week         int
	Production   string
	ProductCode  string
	Product      *Product
	Model        *ModelInfo
	Reason       string
}

func Decode(serial string) Result {
	serial = strings.ToUpper(strings.TrimSpace(serial))

	result := Result{
		Serial: serial,
	}

	switch len(serial) {
	case 11:
		return decode11(serial, result)

	case 12:
		return decode12(serial, result)

	case 10:
		result.Format = FormatRandomized
		result.Valid = false
		result.Reason = "Apple randomized serial format is outside LegacyInspector scope."
		return result

	default:
		result.Format = FormatInvalid
		result.Valid = false
		result.Reason = "Unsupported serial length."
		return result
	}
}
