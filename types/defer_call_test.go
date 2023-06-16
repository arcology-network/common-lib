package types

import (
	"reflect"
	"testing"
)

func TestDeferCallEncoding(t *testing.T) {
	dc := &DeferredCall{
		DeferID:         "123",
		ContractAddress: "45678",
		Signature:       "xxxx",
	}

	buffer := make([]byte, dc.Size())
	dc.EncodeToBuffer(buffer)

	out := (&DeferredCall{}).Decode(buffer)
	if !reflect.DeepEqual(dc, out) {
		t.Error("Error")
	}
}
