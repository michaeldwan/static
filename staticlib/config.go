package staticlib

import (
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/michaeldwan/static/printer"
	"github.com/spf13/cast"
)

const (
	ConfigFileName  string = "static.yml"
	s3RegionDefault string = "us-east-1"
)

type Config struct {
	Path            string
	SourceDirectory string
	S3Region        string
	S3Bucket        string
	redirects       map[string]string
	ignorePatterns  globlist
	gzipPatterns    globlist
	maxAgePatterns  globlist
}

func NewConfig(path string) Config {
	cfg := Config{}
	cfg.loadFromPath(path)
	return cfg
}

func (c *Config) loadFromPath(path string) {
	path, _ = filepath.Abs(path)
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		printer.Infof("Config file %s not found\n", path)
		os.Exit(1)
	}
	c.load(path, file)
}

func (c *Config) load(path string, reader io.Reader) {
	seq, err := parseConfig(reader)
	if err != nil {
		panic(err)
	}
	c.Path = path
	c.loadSourceDirectory(path, seq)
	c.loadS3Region(seq)
	c.loadS3Bucket(seq)
	c.loadRedirects(seq)
	c.loadIgnore(seq)
	if err := c.loadGzip(seq); err != nil {
		panic(err)
	}
	c.loadMaxAge(seq)
}

func (c *Config) loadSourceDirectory(rootDir string, seq sequence) {
	value := cast.ToString(seq.valForKey("source_directory"))
	if value == "" {
		panic("Configuration error: source_directory is missing")
	}
	if !strings.HasPrefix(value, "/") {
		value = filepath.Join(filepath.Dir(rootDir), value)
	}
	c.SourceDirectory = value
}

func (c *Config) loadS3Region(seq sequence) {
	val := cast.ToString(seq.valForKey("s3_region"))
	if val == "" {
		val = s3RegionDefault
	}
	c.S3Region = val
}

func (c *Config) loadS3Bucket(seq sequence) {
	val := cast.ToString(seq.valForKey("s3_bucket"))
	if val == "" {
		panic("Configuration error: s3_bucket is missing")
	}
	c.S3Bucket = val
}

func (c *Config) loadRedirects(seq sequence) {
	// if seq.valForKey("redirects") != "" {
	// 	panic("Configuration error: redirects should be a list of paths and destinations")
	// }
	// TODO: assert right content types (not string val, not value sequence)
	c.redirects = make(map[string]string)
	for _, el := range seq.seqForKey("redirects") {
		path := cast.ToString(el.key)
		dest := cast.ToString(el.value)
		if strings.HasPrefix(path, "/") {
			c.redirects[path[1:]] = dest
		} else {
			c.redirects[path] = dest
		}
	}
}

func (c *Config) Redirects() map[string]string {
	return c.redirects
}

func (c *Config) loadIgnore(seq sequence) {
	c.ignorePatterns = newGlobList(false)
	if el, ok := seq.elForKey("ignore"); ok {
		c.ignorePatterns.loadFromStringSlice(el.stringSliceForSeqValues(), true)
	}

	// Add default ignore patterns
	// TODO: add dot prefix pattern
	// TODO: default windows patterns?
	c.ignorePatterns.add(".DS_Store", true)
	c.ignorePatterns.add(".git", true)
	c.ignorePatterns.add(".svn", true)
}

func (c *Config) ShouldIgnore(path string) bool {
	return cast.ToBool(c.ignorePatterns.get(path))
}

var defaultCompressableMimeTypes = []string{
	"text/html",
	"text/plain",
	"text/css",
	"application/javascript",
	"application/x-javascript",
	"text/xml",
	"application/xml",
	"application/atom+xml",
}

func (c *Config) loadGzip(seq sequence) error {
	c.gzipPatterns = newGlobList(false)
	el, ok := seq.elForKey("gzip")
	if !ok {
		return nil
	}
	if cast.ToBool(el.value) {
		// el.value == `true`, load defaults
		return c.loadGzipValues(defaultCompressableMimeTypes)
	} else if list := el.stringSliceForSeqValues(); len(list) > 0 {
		// el.value == string sequence
		return c.loadGzipValues(list)
	} else if len(el.value) > 0 {
		// el.value == string
		return c.loadGzipValues([]string{el.value})
	}
	return nil
}

func (c *Config) loadGzipValues(values []string) error {
	for _, input := range values {
		mediaType, err := gzipInputToMediaType(input)
		if err != nil {
			return err
		}
		c.gzipPatterns.add(mediaType, true)
	}
	return nil
}

// gzipInputToMediaType converts a string value found in the gzip config sequence
// into a media type string. It accepts either a media type (eg `text/html` as
// defined in RFC 2183), a media type with a wildcard `text/*`) or a file extension.
func gzipInputToMediaType(val string) (string, error) {
	if strings.Contains(val, "/") {
		mediaType, _, err := mime.ParseMediaType(val)
		if err != nil {
			return "", err
		}
		return mediaType, nil
	}
	if !strings.HasPrefix(val, ".") {
		val = fmt.Sprintf(".%s", val)
	}
	if mediaType := mime.TypeByExtension(val); len(mediaType) > 0 {
		return gzipInputToMediaType(mediaType)
	}
	return "", fmt.Errorf("Invalid gzip value: %v", val)
}

func (c *Config) ShouldGzip(contentType string) bool {
	mediaType, _, _ := mime.ParseMediaType(contentType)
	if mediaType != "" {
		return cast.ToBool(c.gzipPatterns.get(mediaType))
	}
	return false
}

func (c *Config) loadMaxAge(seq sequence) {
	c.maxAgePatterns = newGlobList(0)
	el, ok := seq.elForKey("max_age")
	if !ok {
		return
	}
	if defaultMaxAge := cast.ToInt(el.value); defaultMaxAge > 0 {
		c.maxAgePatterns.defaultValue = defaultMaxAge
	} else if len(el.sequence) > 0 {
		for _, subEl := range el.sequence {
			if age := cast.ToInt(subEl.value); age > 0 {
				c.maxAgePatterns.add(subEl.key, age)
			}
		}
	}
}

func (c *Config) MaxAge(path string) int {
	return cast.ToInt(c.maxAgePatterns.get(path))
}
