package tools_test

import (
	"reflect"
	"sokker-org-auto-bidder/tools"
	"testing"
	"time"
)

func TestTimeInZoneBadTimezone(t *testing.T) {
	_, err := tools.TimeInZone("", "now", "bad timezone")
	if err == nil {
		t.Errorf("expected 'ErrBadTimezone' error, got nil")
	}

	switch err.(type) {
	case *tools.ErrBadTimezone:
		return
	default:
		t.Errorf("expected 'ErrBadTimezone' error, got '%s'", reflect.ValueOf(err).Type())
	}
}

func TestTimeInZoneParseError(t *testing.T) {
	_, err := tools.TimeInZone("", "im not time", "Europe/Warsaw")
	if err == nil {
		t.Errorf("expected time parse error, got nil")
	}
}

func TestTimeInZoneSuccess(t *testing.T) {
	timezone := "Etc/GMT+1"
	dt, err := tools.TimeInZone(time.Kitchen, "3:04PM", timezone)
	if err != nil {
		t.Errorf("expcted parsed time, err found: %v", err)
	}

	dtTimezone := dt.Location().String()
	if dtTimezone != timezone {
		t.Errorf("timezone for parsed time equals %s, expected %s", dtTimezone, timezone)
	}

	dtString := dt.In(time.UTC).String()
	expectedUTCString := "0000-01-01 16:04:00 +0000 UTC"
	if dtString != expectedUTCString {
		t.Errorf("time string equals %s, expected %s", dtString, expectedUTCString)
	}
}
