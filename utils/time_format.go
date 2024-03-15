package utils

import (
	"time"
)

func FormatTime(t string) string {
	dob, err := time.Parse("2006-01-02", t)
	if err != nil {
		return "Invalid date format"
	}
	t = dob.Format("2006-01-02")

	return t
}
