package util

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/jbenet/go-base58"
	"io"
)

// Base64
func Base64RawURL(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// B58
func BytesToB58(data []byte) string {
	return base58.Encode(data)
}

func BytesFromB58(b58 string) []byte {
	return base58.Decode(b58)
}

// Hex
func BytesToHex(data []byte) string {
	return hex.EncodeToString(data)
}

func BytesFromHex(hexstr string) []byte {
	data, err := hex.DecodeString(hexstr)
	Check(err)
	return data
}

// JSON
func ToJSON(v interface{}) []byte {
	data, err := json.Marshal(v)
	Check(err)
	return data
}

func FromJSON(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	Check(err)
}

func ReadJSON(r io.Reader, v interface{}) {
	enc := json.NewDecoder(r)
	err := enc.Decode(v)
	Check(err)
}
