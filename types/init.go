/*
 *   Copyright (c) 2024 Arcology Network

 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.

 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.

 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package types

import (
	"encoding/gob"
	"math/big"

	ethCommon "github.com/ethereum/go-ethereum/common"
)

func init() {
	gob.Register(&InclusiveList{})
	gob.Register(&ReapingList{})
	gob.Register(&ReceiptHashList{})

	gob.Register(&StandardTransaction{})

	gob.Register([]*StandardTransaction{})
	gob.Register([]*StandardTransaction{})

	gob.Register([][]byte{})
	gob.Register([]byte{})

	gob.Register(&big.Int{})

	gob.Register(map[ethCommon.Hash]ethCommon.Hash{})

	gob.Register(&IncomingTxs{})

}
