package util

import (
	"io"
	"io/ioutil"
)

func ReadAll(r io.Reader) []byte {
	bytes, err := ioutil.ReadAll(r)
	Check(err)
	return bytes
}
