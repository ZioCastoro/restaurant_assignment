package payloads

import (
	"time"
)

// parseTime parses a pointer to a time string to a time.Time.
// Errors are ignored because it is intended to be used on validated data.
func parseTime(ptr *string) *time.Time {
	if ptr == nil {
		return nil
	}

	res, _ := time.Parse(time.RFC3339, *ptr)

	return &res
}
