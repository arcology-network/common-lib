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

package scheduletype

import "github.com/arcology-network/common-lib/exp/slice"

const (
	SHORT_CONTRACT_ADDRESS_LENGTH = 8 //8 bytes for address
	FUNCTION_SIGNATURE_LENGTH     = 4 // 4 bytes for signature
	CALLEE_ID_LENGTH              = SHORT_CONTRACT_ADDRESS_LENGTH + FUNCTION_SIGNATURE_LENGTH
	MAX_CONFLICT_RATIO            = 0.5
	MAX_NUM_CONFLICTS             = 256

	PROPERTY_PATH        = "func/"
	PROPERTY_PATH_LENGTH = len(PROPERTY_PATH)
	EXECUTION_METHOD     = "execution"
	EXECUTION_EXCEPTED   = "except/"
	DEFERRED_FUNC        = "defer"

	PARALLEL_EXECUTION   = uint8(0) // The default method
	SEQUENTIAL_EXECUTION = uint8(255)
)

// Get the callee key from a message
// func ToKey(msg *eucommon.StandardMessage) string {
// 	if (*msg.Native).To == nil {
// 		return ""
// 	}

// 	if len(msg.Native.Data) == 0 {
// 		return string((*msg.Native.To)[:schtype.FUNCTION_SIGNATURE_LENGTH])
// 	}
// 	return CallToKey((*msg.Native.To)[:], msg.Native.Data[:schtype.FUNCTION_SIGNATURE_LENGTH])
// }

func CallToKey(addr []byte, funSign []byte) string {
	return string(addr[:FUNCTION_SIGNATURE_LENGTH]) + string(funSign[:FUNCTION_SIGNATURE_LENGTH])
}

// The function creates a compact representation of the callee information
func Compact(addr []byte, funSign []byte) []byte {
	addr = slice.Clone(addr) // Make sure the original data is not modified
	return append(addr[:SHORT_CONTRACT_ADDRESS_LENGTH], funSign[:FUNCTION_SIGNATURE_LENGTH]...)
}
