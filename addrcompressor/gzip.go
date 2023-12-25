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
