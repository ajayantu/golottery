package helpers

import (
	"fmt"
	"time"
)

func CompareDates(stringDate string, timeDate time.Time) (bool, error) {
	// Parse the string date (format: dd/mm/yy)
	layout := "02/01/2006"
	parsedDate, err := time.Parse(layout, stringDate)
	if err != nil {
		return false, fmt.Errorf("error parsing string date")
	}

	// Compare the dates
	return parsedDate.After(timeDate), nil
}
