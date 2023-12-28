package datastore

import (
	"fmt"
	"testing"
	"time"

	common "github.com/arcology-network/common-lib/common"
)

func TestCachePolicy(t *testing.T) {
	t0 := time.Now()
	fmt.Println("CachePolicy FreeMemory:", time.Since(t0))
	values := []interface{}{nil, nil, 1, 2}
	// common.RemoveIf(&values, func(v interface{}) bool { return v == nil })
	common.RemoveIf(&values, func(v interface{}) bool { return v == nil })
	if len(values) != 2 || values[0] != 1 || values[1] != 2 {
		t.Error("Error: Expected [1, 2] actual: ", values)
	}

	values = []interface{}{1, nil, nil, 2}
	common.RemoveIf(&values, func(v interface{}) bool { return v == nil })
	if len(values) != 2 || values[0] != 1 || values[1] != 2 {
		t.Error("Error: Expected [1, 2], actual: ", values)
	}

	values = []interface{}{1, 2, nil, nil}
	common.RemoveIf(&values, func(v interface{}) bool { return v == nil })
	if len(values) != 2 || values[0] != 1 || values[1] != 2 {
		t.Error("Error: Expected [1, 2], actual: ", values)
	}

	values = []interface{}{1, nil, 2, nil}
	common.RemoveIf(&values, func(v interface{}) bool { return v == nil })
	if len(values) != 2 || values[0] != 1 || values[1] != 2 {
		t.Error("Error: Expected [1, 2], actual: ", values)
	}

	values = []interface{}{1, 2}
	common.RemoveIf(&values, func(v interface{}) bool { return v == nil })
	if len(values) != 2 || values[0] != 1 || values[1] != 2 {
		t.Error("Error: Expected [1, 2], actual: ", values)
	}

	values = []interface{}{nil, nil}
	common.RemoveIf(&values, func(v interface{}) bool { return v == nil })
	if len(values) != 0 {
		t.Error("Error: Expected [], actual: ", values)
	}
}
