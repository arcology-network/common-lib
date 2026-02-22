package types

import (
	"bytes"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
)

func TestEncodeExecutionResultRoundTrip(t *testing.T) {
	original := &core.ExecutionResult{
		UsedGas:         21000,
		RefundedGas:     1000,
		ReturnData:      []byte{0x01, 0x02, 0x03},
		ContractAddress: common.HexToAddress("0x2222222222222222222222222222222222222222"),
	}

	payload, err := EncodeExecutionResult(original)
	if err != nil {
		t.Fatalf("EncodeExecutionResult returned error: %v", err)
	}

	decoded, err := DecodeExecutionResult(payload)
	if err != nil {
		t.Fatalf("DecodeExecutionResult returned error: %v", err)
	}

	if decoded.UsedGas != original.UsedGas {
		t.Fatalf("UsedGas mismatch: got %d want %d", decoded.UsedGas, original.UsedGas)
	}
	if decoded.RefundedGas != original.RefundedGas {
		t.Fatalf("RefundedGas mismatch: got %d want %d", decoded.RefundedGas, original.RefundedGas)
	}
	if !bytes.Equal(decoded.ReturnData, original.ReturnData) {
		t.Fatalf("ReturnData mismatch: got %x want %x", decoded.ReturnData, original.ReturnData)
	}
	if decoded.ContractAddress != original.ContractAddress {
		t.Fatalf("ContractAddress mismatch: got %s want %s", decoded.ContractAddress.Hex(), original.ContractAddress.Hex())
	}
	if decoded.Err != nil {
		t.Fatalf("Err field expected nil, got %v", decoded.Err)
	}
}

func TestEncodeExecutionResultWithErrorField(t *testing.T) {
	original := &core.ExecutionResult{Err: errors.New("encode failure")}

	payload, err := EncodeExecutionResult(original)
	if err != nil {
		t.Fatalf("EncodeExecutionResult returned error: %v", err)
	}

	decoded, err := DecodeExecutionResult(payload)
	if err != nil {
		t.Fatalf("DecodeExecutionResult returned error: %v", err)
	}
	if decoded.Err == nil {
		t.Fatalf("Err field expected non-nil")
	}
	if decoded.Err.Error() != original.Err.Error() {
		t.Fatalf("Err mismatch: got %q want %q", decoded.Err.Error(), original.Err.Error())
	}
}
