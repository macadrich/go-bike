package utils

import "time"

func ValidateTimestamp(t string) bool {
	layout := time.RFC3339
	_, err := time.Parse(layout, t)

	return err == nil
}
