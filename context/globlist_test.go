package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPatternRegexp(t *testing.T) {
	tests := map[string]string{
		`*`:                       `.{0,}`,
		`*.html`:                  `.{0,}\.html`,
		`index.*`:                 `index\..{0,}`,
		`directory/*`:             `directory\/.{0,}`,
		`*/directory/index.html`:  `.{0,}\/directory\/index\.html`,
		`directory/*/index.html`:  `directory\/[^\/]+\/index\.html`,
		`directory/**/index.html`: `directory\/.{0,}\/index\.html`,
		`index.htm?`:              `index\.htm[^\/]?`,
		`/directory/index.html`:   `^\/directory\/index\.html`,
		`/*.html`:                 `^\/[^\/]{0,}\.html`,
	}

	for input, expected := range tests {
		actual := globRegexpFromPattern(input)
		assert.Equal(t, expected, actual.String())
	}
}

func TestBoolForPath(t *testing.T) {
	tests := map[string]bool{
		`index.html`:        true,
		`abc/index.html`:    true,
		`index.htm`:         false,
		`/assets/image.png`: true,
		`/image.png`:        false,
		`styles.css`:        true,
	}

	g := newGlobList(false)
	g.loadFromMapStringInterface(map[string]interface{}{
		`*.html`:  true,
		`*.css`:   true,
		`/assets`: true,
	})

	for input, expected := range tests {
		actual := g.get(input)
		assert.Equal(t, expected, actual)
	}
}

func TestGlobListGetDefault(t *testing.T) {
	g := newGlobList("hello")
	assert.Equal(t, "hello", g.get("missing"))
}
