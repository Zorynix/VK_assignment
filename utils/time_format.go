package utils

import (
	"time"

	"github.com/rs/zerolog/log"
)

// FormatTime takes a date string 't' as input and attempts to parse it according to a specific layout ("2006-01-02"),
// which is the ISO 8601 format (commonly used international date format).
//
// If the date string is successfully parsed, it reformats the date back into the same ISO 8601 format and returns it.
// This might seem redundant since it returns the date in the same format as the input,
// but it ensures that the input date string is in the correct and expected format. It can also be easily modified
// to return the date in a different format if needed.
//
// In case of a parsing error, it logs the error along with the problematic date string for debugging purposes,
// and returns a predefined error message "Invalid date format". This indicates that the input string
// did not conform to the expected date format.
//
// This function is useful for validating and standardizing date strings in your application,
// ensuring they adhere to a consistent format before they are used in database operations, calculations,
// or displayed in the UI.
//
// Parameters:
// - t (string): The date string to parse and format.
//
// Returns:
// - string: The formatted date string if successful, or an error message if parsing fails.
func FormatTime(t string) string {
	dob, err := time.Parse("2006-01-02", t)
	if err != nil {
		log.Error().Err(err).Str("date", t).Msg("Invalid date format in FormatTime")
		return "Invalid date format"
	}
	t = dob.Format("2006-01-02")

	return t
}
