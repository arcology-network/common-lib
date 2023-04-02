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

func TestSum(t *testing.T) {
	v := []uint{1, 2, 3, 4}
	if Sum(v, uint(0)) != 10 {
		t.Error("Error: Should be 10")
	}

	bytes := []byte{1, 2, 3, 4}
	if Sum(bytes, uint(0)) != 10 {
		t.Error("Error: Should be 10")
	}

}
