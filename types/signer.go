package types

import (
	"math/big"

	evmTypes "github.com/ethereum/go-ethereum/core/types"
)

const (
	Signer_London = iota
	Signer_Cancun
)

func GetSigner(signerType uint8, chainId *big.Int) evmTypes.Signer {
	switch signerType {
	case Signer_London:
		return evmTypes.NewLondonSigner(chainId)
	case Signer_Cancun:
		return evmTypes.NewCancunSigner(chainId)
	}
	return nil
}
