package push

import (
	"bytes"
	"os"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
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
	assert.Equal(t, f.Name(), r.path)
	assert.Equal(t, expected, string(actual))
	assert.Equal(t, "text/html", r.ContentType())
	assert.Equal(t, "http://destination.url", r.RedirectUrl())
}
