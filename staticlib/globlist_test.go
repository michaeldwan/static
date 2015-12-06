package staticlib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGloblistPatternRegexp(t *testing.T) {
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

func TestGloblistBoolForPath(t *testing.T) {
	tests := map[string]bool{
		`index.html`:        true,
		`abc/index.html`:    true,
		`index.htm`:         false,
		`/assets/image.png`: true,
		`/image.png`:        false,
		`styles.css`:        true,
	}

	g := newGlobList(false)
	g.add(`*.html`, true)
	g.add(`*.css`, true)
	g.add(`/assets`, true)

	for input, expected := range tests {
		actual := g.get(input)
		assert.Equal(t, expected, actual)
	}
}

func TestGlobListGetDefault(t *testing.T) {
	g := newGlobList("hello")
	assert.Equal(t, "hello", g.get("missing"))
}

func TestGlobListSortOrder(t *testing.T) {
	tests := map[string]int{
		`index.html`:         1,
		`abc/index.html`:     2,
		`abc/xyz/index.html`: 3,
		`something.html`:     4,
	}

	g := newGlobList(false)
	g.add(`abc/index.html`, 2)
	g.add(`abc/*/index.html`, 3)
	g.add(`index.html`, 1)
	g.add(`*`, 4)

	for input, expected := range tests {
		actual := g.get(input)
		assert.Equal(t, expected, actual, input)
	}
}
