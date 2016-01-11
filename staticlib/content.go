package staticlib

import (
	"os"
	"time"
)

type fileRef struct {
	path string
	size int64
}

func newFileRef(path string) (fileRef, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return fileRef{}, err
	}
	return fileRef{path: path, size: fi.Size()}, nil
}

type Content struct {
	Key     string
	ModTime time.Time

	sourceFile fileRef
	tempFile   fileRef

	ContentType     string
	ContentEncoding string
	CacheControl    string
	Digest          []byte
	RedirectUrl     string
	Notes           []string
}

func (c *Content) workingFile() fileRef {
	if c.tempFile.path != "" {
		return c.tempFile
	}
	return c.sourceFile
}

func (c *Content) Size() int64 {
	return c.workingFile().size
}

func (c *Content) IsRedirect() bool { return c.RedirectUrl != "" }

func (c *Content) cleanTempFile() error {
	if c.tempFile.path == "" {
		return nil
	}
	return os.Remove(c.tempFile.path)
}

func (c *Content) setTempFile(f fileRef) error {
	if err := c.cleanTempFile(); err != nil {
		return err
	}
	c.tempFile = f
	return nil
}
