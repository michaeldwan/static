package staticlib

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testConfigCallback func(c Config, seq sequence)

func newTestConfig(data string, callback testConfigCallback) {
	c := Config{}
	seq, err := parseConfig(strings.NewReader(data))
	if err != nil {
		panic(err)
	}
	callback(c, seq)
}

// func TestConfigPath(t *testing.T) {
// 	path := "/Users/md/src/static/static.yml"
// 	newTestConfig("source_directory: .\ns3_bucket: bucket", func(c Config, seq sequence) {
// 		c.load(path, seq)
// 		assert.Equal(t, path, c.Path)
// 	})
// }

func TestConfigSourceDirectory(t *testing.T) {
	newTestConfig("source_directory: ./build", func(c Config, seq sequence) {
		c.loadSourceDirectory("/Users/md/src/static/static.yml", seq)
		assert.Equal(t, "/Users/md/src/static/build", c.SourceDirectory)
	})
}

func TestConfigSourceDirectoryAbsolute(t *testing.T) {
	newTestConfig("source_directory: /src/build", func(c Config, seq sequence) {
		c.loadSourceDirectory("/Users/md/src/static/static.yml", seq)
		assert.Equal(t, "/src/build", c.SourceDirectory)
	})
}

func TestConfigSourceDirectoryMissing(t *testing.T) {
	newTestConfig("source_directory: ", func(c Config, seq sequence) {
		assert.Panics(t, func() {
			c.loadSourceDirectory("/Users/md/src/static/static.yml", seq)
		})
	})

	newTestConfig("", func(c Config, seq sequence) {
		assert.Panics(t, func() {
			c.loadSourceDirectory("/Users/md/src/static/static.yml", seq)
		})
	})
}

func TestConfigS3Region(t *testing.T) {
	newTestConfig("s3_region: us-west-1", func(c Config, seq sequence) {
		c.loadS3Region(seq)
		assert.Equal(t, "us-west-1", c.S3Region)
	})
}

func TestConfigS3RegionMissing(t *testing.T) {
	newTestConfig("s3_region: ", func(c Config, seq sequence) {
		c.loadS3Region(seq)
		assert.Equal(t, s3RegionDefault, c.S3Region)
	})

	newTestConfig("", func(c Config, seq sequence) {
		c.loadS3Region(seq)
		assert.Equal(t, s3RegionDefault, c.S3Region)
	})
}

func TestConfigS3Bucket(t *testing.T) {
	newTestConfig("s3_bucket: michaeldwan.com", func(c Config, seq sequence) {
		c.loadS3Bucket(seq)
		assert.Equal(t, "michaeldwan.com", c.S3Bucket)
	})
}

func TestConfigS3BucketMissing(t *testing.T) {
	newTestConfig("s3_bucket: ", func(c Config, seq sequence) {
		assert.Panics(t, func() {
			c.loadS3Bucket(seq)
		})
	})

	newTestConfig("", func(c Config, seq sequence) {
		assert.Panics(t, func() {
			c.loadS3Bucket(seq)
		})
	})
}

func TestConfigRedirects(t *testing.T) {
	var data = `
redirects:
  left: right
  up: down
  /slash/prefix: noslash
`
	newTestConfig(data, func(c Config, seq sequence) {
		c.loadRedirects(seq)
		assert.Len(t, c.redirects, 3)
		assert.Equal(t, "right", c.Redirects()["left"])
		assert.Equal(t, "down", c.Redirects()["up"])
		assert.Equal(t, "noslash", c.Redirects()["slash/prefix"])
	})
}

func TestConfigRedirectsEmpty(t *testing.T) {
	var data = "redirects:  "
	newTestConfig(data, func(c Config, seq sequence) {
		c.loadRedirects(seq)
		assert.Equal(t, 0, len(c.redirects))
	})
}

func TestConfigRedirectsWrongType(t *testing.T) {
	var data = "redirects: 123"
	newTestConfig(data, func(c Config, seq sequence) {
		assert.Panics(t, func() {
			c.loadRedirects(seq)
		})
		assert.Equal(t, 0, len(c.redirects))
	})
}

func TestConfigIgnore(t *testing.T) {
	var data = `
ignore:
  - "*.html"
  - "assets/*"
`
	newTestConfig(data, func(c Config, seq sequence) {
		c.loadIgnore(seq)
		assert.Equal(t, true, c.ShouldIgnore("index.html"))
		assert.Equal(t, true, c.ShouldIgnore("assets/image.png"))
		assert.Equal(t, false, c.ShouldIgnore("style.css"))
		assert.Equal(t, true, c.ShouldIgnore(".git"))
		assert.Equal(t, true, c.ShouldIgnore("directory/.git"))
		assert.Equal(t, true, c.ShouldIgnore(".DS_Store"))
		assert.Equal(t, true, c.ShouldIgnore("directory/.DS_Store"))
		assert.Equal(t, true, c.ShouldIgnore(".svn"))
		assert.Equal(t, true, c.ShouldIgnore("directory/.svn"))
	})
}

func TestConfigIgnoreEmpty(t *testing.T) {
	var data = `
ignore:
`
	newTestConfig(data, func(c Config, seq sequence) {
		c.loadIgnore(seq)
		assert.Equal(t, false, c.ShouldIgnore("index.html"))
		assert.Equal(t, false, c.ShouldIgnore("assets/image.png"))
		assert.Equal(t, false, c.ShouldIgnore("style.css"))
		assert.Equal(t, true, c.ShouldIgnore(".git"))
		assert.Equal(t, true, c.ShouldIgnore("directory/.git"))
		assert.Equal(t, true, c.ShouldIgnore(".DS_Store"))
		assert.Equal(t, true, c.ShouldIgnore("directory/.DS_Store"))
		assert.Equal(t, true, c.ShouldIgnore(".svn"))
		assert.Equal(t, true, c.ShouldIgnore("directory/.svn"))
	})
}

func TestConfigGzipList(t *testing.T) {
	var data = `
gzip:
  - "*.html"
  - "assets/*"
`
	newTestConfig(data, func(c Config, seq sequence) {
		c.loadGzip(seq)
		assert.Equal(t, true, c.ShouldGzip("index.html"))
		assert.Equal(t, true, c.ShouldGzip("assets/image.png"))
		assert.Equal(t, false, c.ShouldGzip("style.css"))
		assert.Equal(t, true, c.ShouldGzip("assets/style.css"))
	})
}

func TestConfigGzipTrue(t *testing.T) {
	var data = `
gzip: true
`
	newTestConfig(data, func(c Config, seq sequence) {
		c.loadGzip(seq)
		assert.Equal(t, true, c.ShouldGzip("index.html"))
		assert.Equal(t, false, c.ShouldGzip("assets/image.png"))
		assert.Equal(t, true, c.ShouldGzip("style.css"))
		assert.Equal(t, true, c.ShouldGzip("assets/style.css"))
		assert.Equal(t, true, c.ShouldGzip("assets/app.js"))
	})
}

func TestConfigGzipFalse(t *testing.T) {
	var data = `
gzip: false
`
	newTestConfig(data, func(c Config, seq sequence) {
		c.loadGzip(seq)
		assert.Empty(t, c.gzipPatterns.globs)
	})
}

func TestConfigGzipMissing(t *testing.T) {
	var data = ""
	newTestConfig(data, func(c Config, seq sequence) {
		c.loadGzip(seq)
		assert.Empty(t, c.gzipPatterns.globs)
	})
}

func TestConfigMaxAgePatterns(t *testing.T) {
	var data = `
max_age:
  "*.html": 123
  "*.css": 456
`
	newTestConfig(data, func(c Config, seq sequence) {
		c.loadMaxAge(seq)
		assert.Equal(t, 123, c.maxAgePatterns.get("index.html"))
		assert.Equal(t, 456, c.maxAgePatterns.get("style.css"))
		assert.Equal(t, 0, c.maxAgePatterns.get("missing.go"))
	})
}

func TestConfigMaxAgePatternPriority(t *testing.T) {
	var data = `
max_age:
  "abc/index.html": 456
  "index.html": 123
  "*": 789
`
	newTestConfig(data, func(c Config, seq sequence) {
		c.loadMaxAge(seq)
		assert.Equal(t, 123, c.maxAgePatterns.get("index.html"))
		assert.Equal(t, 456, c.maxAgePatterns.get("abc/index.html"))
		assert.Equal(t, 123, c.maxAgePatterns.get("xyz/index.html"))
		assert.Equal(t, 789, c.maxAgePatterns.get("xyz/styles.css"))
	})
}

func TestConfigMaxAgeDefault(t *testing.T) {
	var data = `
max_age: 123
`
	newTestConfig(data, func(c Config, seq sequence) {
		c.loadMaxAge(seq)
		assert.Equal(t, 123, c.maxAgePatterns.defaultValue)
	})
}

func TestConfigMaxAgeMissing(t *testing.T) {
	var data = ""
	newTestConfig(data, func(c Config, seq sequence) {
		c.loadMaxAge(seq)
		assert.Equal(t, 0, c.maxAgePatterns.defaultValue)
	})
}
