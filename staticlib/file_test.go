package staticlib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileContentType(t *testing.T) {
	f := newFile("file/key.json", "path/to/file.json")
	f.contentType = "text/json"
	assert.Equal(t, "text/json", f.ContentType())
}

func TestFileContentTypeDefault(t *testing.T) {
	f := newFile("file/key.json", "path/to/file.json")
	f.contentType = ""
	assert.Equal(t, "binary/octet-stream", f.ContentType())
}

func TestFileIsRedirect(t *testing.T) {
	f := newFile("file/key.json", "path/to/file.json")
	assert.False(t, f.IsRedirect())
	f.redirectUrl = "http://somewhere.com"
	assert.True(t, f.IsRedirect())
}
