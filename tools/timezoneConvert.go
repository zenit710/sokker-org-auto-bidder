package tools

import (
	"fmt"
	"time"
)

// ErrBadTimezone is raised when provided timezone is unknown
type ErrBadTimezone struct {
	name string
}

func (e *ErrBadTimezone) Error() string {
	return fmt.Sprintf("Timezone '%s' is not known", e.name)
}

// TimeInZone parses date string and returns time.Time in chosen timezone (by string)
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
