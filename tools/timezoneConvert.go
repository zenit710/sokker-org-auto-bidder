package tools

import (
	"fmt"
	"time"
)

type ErrBadTimezone struct {
	name string
}

func (e *ErrBadTimezone) Error() string {
	return fmt.Sprintf("Timezone '%s' is not known", e.name)
}

// Get time in chosen Timezone (by string)
func TimeInZone(layout string, value string, zone string) (time.Time, error) {
	orgLoc, err := time.LoadLocation(zone)
	if err != nil {
		return time.Time{}, &ErrBadTimezone{zone}
	}

	dt, err := time.ParseInLocation(layout, value, orgLoc)
	if err != nil {
		return time.Time{}, err
	}

	return dt, nil
}
