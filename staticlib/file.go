package staticlib

import (
	"fmt"
	"os"
)

type File struct {
	Key             string
	Path            string
	Size            int64
	ContentType     string
	ContentEncoding string
	Digest          []byte
	CacheControl    string
	RedirectUrl     string
	Notes           []string
}

func newFile(key, path string) File {
	return File{Key: key, Path: path}
}

func (f File) IsRedirect() bool { return f.RedirectUrl != "" }
func (f File) Desc() string {
	if f.IsRedirect() {
		return fmt.Sprintf("%s --> %s", f.Key, f.RedirectUrl)
	}
	return fmt.Sprintf("%s", f.Key)
}

func (f File) Body() *os.File {
	out, err := os.Open(f.Path)
	if err != nil {
		panic(err)
	}
	return out
}
