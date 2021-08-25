package signature

import "testing"

func TestGenerator(t *testing.T) {
	sigs := GetParallelFuncList()
	for _, sig := range sigs {
		t.Log(sig)
	}
}
