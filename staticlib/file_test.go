package staticlib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileIsRedirect(t *testing.T) {
	f := newFile("file/key.json", "path/to/file.json")
	assert.False(t, f.IsRedirect())
	f.RedirectUrl = "http://somewhere.com"
	assert.True(t, f.IsRedirect())
}
