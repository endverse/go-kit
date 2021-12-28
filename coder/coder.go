package coder

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
)

var b64 = base64.StdEncoding

var magicGzip = []byte{0x1f, 0x8b, 0x08}

// Encode encodes a object returning a base64 encoded
// gzipped string representation, or error.
func Encode(src interface{}) (string, error) {
	b, err := json.Marshal(src)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	w, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return "", err
	}
	if _, err = w.Write(b); err != nil {
		return "", err
	}
	w.Close()

	return b64.EncodeToString(buf.Bytes()), nil
}

// Decode decodes the bytes of data into a object
// type. Data must contain a base64 encoded gzipped string,
// otherwise an error is returned.
func Decode(data string, dest interface{}) error {
	// base64 decode string
	b, err := b64.DecodeString(data)
	if err != nil {
		return err
	}

	if bytes.Equal(b[0:3], magicGzip) {
		r, err := gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			return err
		}
		defer r.Close()
		b2, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}
		b = b2
	}

	// unmarshal object bytes
	if err := json.Unmarshal(b, dest); err != nil {
		return err
	}
	return nil
}
