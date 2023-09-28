package datacompression

import (
	"bytes"
	"compress/gzip"
	"io"
)

func CompressGZip(buffer []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(buffer); err != nil {
		return []byte{}, err
	}

	if err := w.Close(); err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func DecompressGZip(compressed []byte) ([]byte, error) {
	var buf bytes.Buffer
	buf.Write(compressed)
	reader, err := gzip.NewReader(&buf)
	if err != nil {
		return []byte{}, err
	}

	if err := reader.Close(); err != nil {
		return []byte{}, err
	}
	return io.ReadAll(reader)
}
