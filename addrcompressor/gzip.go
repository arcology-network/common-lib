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

package addrcompressor

import (
	"bytes"
	"compress/gzip"
	"io"
)

func CompressGZip(buffer []byte, name, comment string) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	writer.Name = name
	writer.Comment = comment

	if _, err := writer.Write(buffer); err != nil {
		return []byte{}, err
	}

	if err := writer.Close(); err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func DecompressGZip(compressed []byte) ([]byte, string, string, error) {
	var buf bytes.Buffer
	buf.Write(compressed)
	reader, err := gzip.NewReader(&buf)
	if err != nil {
		return []byte{}, "", "", err
	}

	if err := reader.Close(); err != nil {
		return []byte{}, "", "", err
	}

	uncompressed, err := io.ReadAll(reader)
	if err != nil {
		return []byte{}, "", "", err
	}
	return uncompressed, reader.Name, reader.Comment, err
}
