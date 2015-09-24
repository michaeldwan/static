package push

import (
	"fmt"
	"os"
)

type File struct {
	key             string
	path            string
	size            int64
	contentType     string
	contentEncoding string
	digest          []byte
	cacheControl    string
	redirectUrl     string
}

func newFile(key, path string) *File {
	f := &File{key: key}
	f.setPath(path)
	return f
}

func (f *File) Key() string             { return f.key }
func (f *File) ContentEncoding() string { return f.contentEncoding }
func (f *File) Size() int64             { return f.size }
func (f *File) Digest() []byte          { return f.digest }
func (f *File) CacheControl() string    { return f.cacheControl }
func (f *File) RedirectUrl() string     { return f.redirectUrl }
func (f *File) IsRedirect() bool        { return f.redirectUrl != "" }
func (f *File) Desc() string {
	if f.IsRedirect() {
		return fmt.Sprintf("%s --> %s", f.Key(), f.RedirectUrl())
	}
	return fmt.Sprintf("%s", f.Key())
}

func (f *File) ContentType() string {
	if f.contentType == "" {
		return "binary/octet-stream"
	}
	return f.contentType
}

func (f *File) setPath(path string) {
	f.path = path
}

func (f *File) Body() *os.File {
	out, err := os.Open(f.path)
	if err != nil {
		panic(err)
	}
	return out
}
