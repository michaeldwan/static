package staticlib

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testRedirectPage(destination string) string {
	return strings.Replace(redirBody, "{{.}}", destination, 1)
}

func TestRedirectRender(t *testing.T) {
	expected := testRedirectPage("http://destination.url")
	var writer bytes.Buffer
	renderRedirect("http://destination.url", &writer)
	assert.Equal(t, expected, writer.String())
}

func TestNewRedirect(t *testing.T) {
	expected := testRedirectPage("http://destination.url")
	f, _ := ioutil.TempFile("", "")
	defer os.Remove(f.Name())
	r := newRedirect(f, "redirect", "http://destination.url")
	actual, _ := ioutil.ReadFile(f.Name())
	assert.Equal(t, f.Name(), r.Path)
	assert.Equal(t, expected, string(actual))
	assert.Equal(t, "http://destination.url", r.RedirectUrl)
}
