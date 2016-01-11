package staticlib

import (
	"crypto/md5"
	"io"
)

func digestReader(r io.Reader) ([]byte, error) {
	hash := md5.New()
	var result []byte
	if _, err := io.Copy(hash, r); err != nil {
		return nil, err
	}
	return hash.Sum(result), nil
}
