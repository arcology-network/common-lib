package cachepolicy

import (
	"fmt"
	"testing"
	"time"

	slice "github.com/arcology-network/common-lib/exp/slice"
)

func TestCachePolicy(t *testing.T) {
	t0 := time.Now()
	fmt.Println("CachePolicy FreeMemory:", time.Since(t0))
	values := []interface{}{nil, nil, 1, 2}
	// common.RemoveIf(&values, func(v interface{}) bool { return v == nil })
	slice.RemoveIf(&values, func(_ int, v interface{}) bool { return v == nil })
	if len(values) != 2 || values[0] != 1 || values[1] != 2 {
		t.Error("Error: Expected [1, 2] actual: ", values)
	}

	values = []interface{}{1, nil, nil, 2}
	slice.RemoveIf(&values, func(_ int, v interface{}) bool { return v == nil })
	if len(values) != 2 || values[0] != 1 || values[1] != 2 {
		t.Error("Error: Expected [1, 2], actual: ", values)
	}

	values = []interface{}{1, 2, nil, nil}
	slice.RemoveIf(&values, func(_ int, v interface{}) bool { return v == nil })
	if len(values) != 2 || values[0] != 1 || values[1] != 2 {
		t.Error("Error: Expected [1, 2], actual: ", values)
	}

	values = []interface{}{1, nil, 2, nil}
	slice.RemoveIf(&values, func(_ int, v interface{}) bool { return v == nil })
	if len(values) != 2 || values[0] != 1 || values[1] != 2 {
		t.Error("Error: Expected [1, 2], actual: ", values)
	}

	values = []interface{}{1, 2}
	slice.RemoveIf(&values, func(_ int, v interface{}) bool { return v == nil })
	if len(values) != 2 || values[0] != 1 || values[1] != 2 {
		t.Error("Error: Expected [1, 2], actual: ", values)
	}

	values = []interface{}{nil, nil}
	slice.RemoveIf(&values, func(_ int, v interface{}) bool { return v == nil })
	if len(values) != 0 {
		t.Error("Error: Expected [], actual: ", values)
	}
}
