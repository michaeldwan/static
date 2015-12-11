package staticlib

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

type glob struct {
	pattern *regexp.Regexp
	value   interface{}
}

type globlist struct {
	defaultValue interface{}
	globs        []glob
}

func newGlobList(defaultValue interface{}) globlist {
	return globlist{
		defaultValue: defaultValue,
		globs:        make([]glob, 0),
	}
}

func (g *globlist) loadFromStringSlice(in []string, value interface{}) {
	for _, pattern := range in {
		g.add(pattern, value)
	}
}

func (g *globlist) add(pattern string, value interface{}) {
	r := globRegexpFromPattern(pattern)
	glob := glob{r, value}
	g.globs = append(g.globs, glob)
}

func (g *globlist) get(pattern string) interface{} {
	for _, glob := range g.globs {
		if glob.pattern.MatchString(pattern) {
			return glob.value
		}
	}
	return g.defaultValue
}

var boundaryGlobRegexp = regexp.MustCompile(`^\*|\*$`)
var escapeRegexp = regexp.MustCompile(`([\.\/\+])`)

func globRegexpFromPattern(pattern string) *regexp.Regexp {
	// TODO: use system path separator instead of hardcoded slash
	// TODO: add dot prefix pattern to simplify ignoring hidden files & dirs
	pattern = filepath.Clean(pattern)

	// escape special characters
	pattern = escapeRegexp.ReplaceAllString(pattern, `\$1`)
	// single character wildcard
	pattern = strings.Replace(pattern, `?`, `[^\/]?`, -1)
	// boundary wildcard
	pattern = boundaryGlobRegexp.ReplaceAllString(pattern, `.{0,}`)
	// single directory wildcard
	pattern = strings.Replace(pattern, `\/*\/`, `\/[^\/]+\/`, -1)
	// recursive directory wildcard
	pattern = strings.Replace(pattern, `\/**\/`, `\/.{0,}\/`, -1)
	// directory name or file component wildcard
	pattern = strings.Replace(pattern, `*`, `[^\/]{0,}`, -1)
	// match start of path if leading with /
	if strings.HasPrefix(pattern, `\/`) {
		pattern = fmt.Sprintf("^%s", pattern)
	}

	r, err := regexp.Compile(pattern)

	if err != nil {
		panic(err)
	}
	return r
}
