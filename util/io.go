package util

import (
	"io"
	"io/ioutil"
)

func Copy(w io.Writer, r io.Reader) {
	_, err := io.Copy(w, r)
	Check(err)
}

func ReadAll(r io.Reader) []byte {
	bytes, err := ioutil.ReadAll(r)
	Check(err)
	return bytes
}
