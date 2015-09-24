package context

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
)

type workingDir struct {
	path string
}

func newWorkingDir(name string) workingDir {
	digest := md5.Sum([]byte(name))
	hexDigest := hex.EncodeToString(digest[:])
	wd := workingDir{}
	wd.path = fmt.Sprintf("/tmp/webmaster/%s", hexDigest)
	wd.clean()
	wd.make()
	return wd
}

func (w workingDir) tempFile() *os.File {
	f, err := ioutil.TempFile(w.path, "")
	if err != nil {
		panic(err)
	}
	return f
}

func (w workingDir) make() {
	if err := os.MkdirAll(w.path, 0755); err != nil {
		panic(err)
	}
}

func (w workingDir) clean() {
	if err := os.RemoveAll(w.path); err != nil {
		panic(err)
	}
}
