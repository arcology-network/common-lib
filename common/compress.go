package common

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"io"
	"io/ioutil"
	"log"
)

//with zlib  compress
func ZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

//with zlib uncompress
func ZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

//with gzip  compress
func GzipCompress(src []byte) []byte {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()
	w.Write(src)
	w.Flush()
	return b.Bytes()
}

//with gzip uncompress
func GzipUnCompress(compressSrc []byte) []byte {
	var b bytes.Buffer

	b.Write(compressSrc)

	r, _ := gzip.NewReader(&b)
	defer r.Close()
	undatas, _ := ioutil.ReadAll(r)

	return undatas
}

//with deflate  compress
func DeflateCompress(src []byte) []byte {

	buf := bytes.NewBuffer(nil)

	flateWrite, err := flate.NewWriter(buf, flate.BestCompression)
	if err != nil {
		log.Fatalln(err)
	}
	defer flateWrite.Close()

	flateWrite.Write(src)
	flateWrite.Flush()
	return buf.Bytes()

}

//with deflate uncompress
func DeflateUnCompress(compressSrc []byte) []byte {
	buf := bytes.NewBuffer(nil)

	_, err := buf.Write(compressSrc)
	if err != nil {
		return nil
	}
	r := flate.NewReader(buf)
	defer r.Close()
	undatas, _ := ioutil.ReadAll(r)

	return undatas
}

func Compress(data []byte) []byte {
	buf := bytes.NewBuffer(nil)
	w := lzw.NewWriter(buf, lzw.LSB, 8)
	w.Write(data)
	w.Close()
	return buf.Bytes()
}

func Decompress(data []byte) []byte {
	buf := bytes.NewBuffer(nil)
	buf.Write(data)
	r := lzw.NewReader(buf, lzw.LSB, 8)
	defer r.Close()
	dest := bytes.NewBuffer(nil)
	io.Copy(dest, r)
	return dest.Bytes()
}
