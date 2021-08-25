package signature

import "github.com/HPISTechnologies/evm/crypto"

func GetParallelFuncList() [][]byte {
	funcs := []string{
		// KittyOwnership
		"balanceOf(address)",
		"transfer(address,uint256)",
		"approve(address,uint256)",
		"transferFrom(address,address,uint256)",
		"ownerOf(uint256)",
		// KittyBreeding
		"approveSiring(address,uint256)",
		"isReadyToBreed(uint256)",
		"canBreedWith(uint256,uint256)",
		"breedWith(uint256,uint256)",
		"breedWithAuto(uint256,uint256)",
		"giveBirth(uint256)",
		// KittyAuction
		"createSaleAuction(uint256,uint256,uint256,uint256)",
		"createSiringAuction(uint256,uint256,uint256,uint256)",
		"bidOnSiringAuction(uint256,uint256)",
		// KittyCore
		"getKitty(uint256)",
		// ClockAuction
		"createAuction(uint256,uint256,uint256,uint256,address)",
		"bid(uint256)",
		"cancelAuction(uint256)",
		"getAuction(uint256)",
		"getCurrentPrice(uint256)",
		//test
		"parallel()",
		"func(uint256,uint256,bytes,address,uint256,address,uint256)",
		//"func(uint256,uint256,address,uint256,address,uint256)",
		"createPromoKitty(uint256,address)",
		//"createGen0Auction(uint256)",
		"increase(uint256)",

		"approve(address)",
		"approve(address,uint256)",
		"transfer(address,uint256)",
		"transferFrom(address,address,uint256)",
		"push(address,uint256)",
		"pull(address,uint256)",
		"move(address,address,uint256)",
		"mint(uint256)",
		"burn(uint256)",
		"mint(address,uint256)",
		"burn(address,uint256)",
	}

	signatures := make([][]byte, 0, len(funcs))
	for _, f := range funcs {
		signatures = append(signatures, crypto.Keccak256([]byte(f))[:4])
	}
	return signatures
}

func GetParallelFuncMap() map[string]int {

	concurrencyTable := make(map[string]int)
	for _, v := range GetParallelFuncList() {
		concurrencyTable[string(v[:4])] = 2
	}
	return concurrencyTable
}
