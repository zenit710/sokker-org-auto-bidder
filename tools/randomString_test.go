package tools_test

import (
	"sokker-org-auto-bidder/tools"
	"testing"
)

func TestString(t *testing.T) {
	for i := 0; i < 10; i++ {
		if strlen := len(tools.String(i)); strlen != i {
			t.Errorf("Generated string length %d not equal to asked %d", strlen, i)
		}
	}
}
