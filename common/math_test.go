package common

import (
	"testing"
)

func TestMinMax(t *testing.T) {
	if Min(1, 9) != 1 {
		t.Error("Error: Should be 1")
	}

	if Max(1, 9) != 9 {
		t.Error("Error: Should be 9")
	}
}
