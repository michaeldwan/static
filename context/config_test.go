package context

import (
	"github.com/go-yaml/yaml"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

type newTestConfigCallback func(c *Config, d configMap)

func newTestConfig(data string, callback newTestConfigCallback) {
	c := newConfig()
	d := make(configMap)
	if err := yaml.Unmarshal([]byte(data), &d); err != nil {
		log.Panic(err)
	}
	callback(c, d)
}

func TestPath(t *testing.T) {
	path := "/Users/md/src/webmaster/webmaster.yml"
	newTestConfig("source_directory: .\ns3_bucket: bucket", func(c *Config, d configMap) {
		c.load(path, d)
		assert.Equal(t, path, c.Path)
	})
}

func TestSourceDirectory(t *testing.T) {
	newTestConfig("source_directory: ./build", func(c *Config, d configMap) {
		c.loadSourceDirectory("/Users/md/src/webmaster/webmaster.yml", d)
		assert.Equal(t, "/Users/md/src/webmaster/build", c.SourceDirectory)
	})
}

func TestSourceDirectoryAbsolute(t *testing.T) {
	newTestConfig("source_directory: /src/build", func(c *Config, d configMap) {
		c.loadSourceDirectory("/Users/md/src/webmaster/webmaster.yml", d)
		assert.Equal(t, "/src/build", c.SourceDirectory)
	})
}

func TestSourceDirectoryMissing(t *testing.T) {
	newTestConfig("source_directory: ", func(c *Config, d configMap) {
		assert.Panics(t, func() {
			c.loadSourceDirectory("/Users/md/src/webmaster/webmaster.yml", d)
		})
	})

	newTestConfig("", func(c *Config, d configMap) {
		assert.Panics(t, func() {
			c.loadSourceDirectory("/Users/md/src/webmaster/webmaster.yml", d)
		})
	})
}

func TestS3Region(t *testing.T) {
	newTestConfig("s3_region: us-west-1", func(c *Config, d configMap) {
		c.loadS3Region(d)
		assert.Equal(t, "us-west-1", c.S3Region)
	})
}

func TestS3RegionMissing(t *testing.T) {
	newTestConfig("s3_region: ", func(c *Config, d configMap) {
		c.loadS3Region(d)
		assert.Equal(t, S3RegionDefault, c.S3Region)
	})

	newTestConfig("", func(c *Config, d configMap) {
		c.loadS3Region(d)
		assert.Equal(t, S3RegionDefault, c.S3Region)
	})
}

func TestS3Bucket(t *testing.T) {
	newTestConfig("s3_bucket: michaeldwan.com", func(c *Config, d configMap) {
		c.loadS3Bucket(d)
		assert.Equal(t, "michaeldwan.com", c.S3Bucket)
	})
}

func TestS3BucketMissing(t *testing.T) {
	newTestConfig("s3_bucket: ", func(c *Config, d configMap) {
		assert.Panics(t, func() {
			c.loadS3Bucket(d)
		})
	})

	newTestConfig("", func(c *Config, d configMap) {
		assert.Panics(t, func() {
			c.loadS3Bucket(d)
		})
	})
}

func TestRedirects(t *testing.T) {
	var data = `
redirects:
  left: right
  up: down
  /slash/prefix: noslash
`
	newTestConfig(data, func(c *Config, d configMap) {
		c.loadRedirects(d)
		assert.Len(t, c.redirects, 3)
		assert.Equal(t, "right", c.Redirects()["left"])
		assert.Equal(t, "down", c.Redirects()["up"])
		assert.Equal(t, "noslash", c.Redirects()["slash/prefix"])
	})
}

func TestRedirectsEmpty(t *testing.T) {
	var data = "redirects:  "
	newTestConfig(data, func(c *Config, d configMap) {
		c.loadRedirects(d)
		assert.Equal(t, 0, len(c.redirects))
	})
}

func TestRedirectsWrongType(t *testing.T) {
	var data = "redirects: 123"
	newTestConfig(data, func(c *Config, d configMap) {
		assert.Panics(t, func() {
			c.loadRedirects(d)
		})
		assert.Equal(t, 0, len(c.redirects))
	})
}

func TestIgnore(t *testing.T) {
	var data = `
ignore:
  - "*.html"
  - "assets/*"
`
	newTestConfig(data, func(c *Config, d configMap) {
		c.loadIgnore(d)
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

func TestIgnoreEmpty(t *testing.T) {
	var data = `
ignore:
`
	newTestConfig(data, func(c *Config, d configMap) {
		c.loadIgnore(d)
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

func TestGzipList(t *testing.T) {
	var data = `
gzip:
  - "*.html"
  - "assets/*"
`
	newTestConfig(data, func(c *Config, d configMap) {
		c.loadGzip(d)
		assert.Equal(t, true, c.ShouldGzip("index.html"))
		assert.Equal(t, true, c.ShouldGzip("assets/image.png"))
		assert.Equal(t, false, c.ShouldGzip("style.css"))
		assert.Equal(t, true, c.ShouldGzip("assets/style.css"))
	})
}

func TestGzipTrue(t *testing.T) {
	var data = `
gzip: true
`
	newTestConfig(data, func(c *Config, d configMap) {
		c.loadGzip(d)
		assert.Equal(t, true, c.ShouldGzip("index.html"))
		assert.Equal(t, false, c.ShouldGzip("assets/image.png"))
		assert.Equal(t, true, c.ShouldGzip("style.css"))
		assert.Equal(t, true, c.ShouldGzip("assets/style.css"))
		assert.Equal(t, true, c.ShouldGzip("assets/app.js"))
	})
}

func TestGzipFalse(t *testing.T) {
	var data = `
gzip: false
`
	newTestConfig(data, func(c *Config, d configMap) {
		c.loadGzip(d)
		assert.Empty(t, c.gzipPatterns.globs)
	})
}

func TestGzipMissing(t *testing.T) {
	var data = ""
	newTestConfig(data, func(c *Config, d configMap) {
		c.loadGzip(d)
		assert.Empty(t, c.gzipPatterns.globs)
	})
}

func TestConfigMaxAgePatterns(t *testing.T) {
	var data = `
max-age:
  "*.html": 123
  "*.css": 456
`
	newTestConfig(data, func(c *Config, d configMap) {
		c.loadMaxAge(d)
		assert.Equal(t, 123, c.maxAgePatterns.get("index.html"))
		assert.Equal(t, 456, c.maxAgePatterns.get("style.css"))
		assert.Equal(t, 0, c.maxAgePatterns.get("missing.go"))
	})
}

func TestConfigMaxAgeDefault(t *testing.T) {
	var data = `
max-age: 123
`
	newTestConfig(data, func(c *Config, d configMap) {
		c.loadMaxAge(d)
		assert.Equal(t, 123, c.maxAgePatterns.defaultValue)
	})
}

func TestConfigMaxAgeMissing(t *testing.T) {
	var data = ""
	newTestConfig(data, func(c *Config, d configMap) {
		c.loadMaxAge(d)
		assert.Equal(t, 0, c.maxAgePatterns.defaultValue)
	})
}
