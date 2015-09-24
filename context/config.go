package context

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/spf13/cast"
)

const ConfigFileName string = "webmaster.yml"

type configMissingError struct {
	err string
}

func (e configMissingError) Error() string { return e.err }
func newConfigMissingError(keyPath string) error {
	return configMissingError{fmt.Sprintf("%s is missing", keyPath)}
}

type configMap map[interface{}]interface{}

func (c *configMap) getE(keyPath string) (interface{}, error) {
	keys := strings.Split(keyPath, ":")
	kv := *c
	for index, key := range keys {
		if val, ok := kv[key]; ok {
			if index == len(keys)-1 {
				return val, nil
			}
			if kv, ok = val.(map[interface{}]interface{}); !ok {
				break
			}
		}
	}
	return nil, newConfigMissingError(keyPath)
}

func (c *configMap) get(keyPath string) interface{} {
	val, _ := c.getE(keyPath)
	return val
}

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

func newConfig() *Config {
	return &Config{
		ignorePatterns: newGlobList(false),
		gzipPatterns:   newGlobList(false),
		maxAgePatterns: newGlobList(0),
		redirects:      make(map[string]string),
	}
}

func (c *Config) loadFile(path string) {
	path, _ = filepath.Abs(path)
	data, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		panic(fmt.Sprintf("Config file %s not found\n", path))
	}
	c.loadPathAndData(path, data)
}

func (c *Config) loadPathAndData(path string, data []byte) {
	raw := make(configMap)
	if err := yaml.Unmarshal(data, &raw); err != nil {
		log.Panic(err)
	}
	c.load(path, raw)
}

func (c *Config) load(path string, raw configMap) {
	c.Path = path
	c.loadSourceDirectory(path, raw)
	c.loadS3Region(raw)
	c.loadS3Bucket(raw)
	c.loadRedirects(raw)
	c.loadIgnore(raw)
	c.loadGzip(raw)
	c.loadMaxAge(raw)
}

func (c *Config) loadSourceDirectory(rootDir string, raw configMap) {
	value := cast.ToString(raw.get("source_directory"))
	if value == "" {
		log.Panicln("Configuration error: source_directory is missing")
	}
	if !strings.HasPrefix(value, "/") {
		value = filepath.Join(filepath.Dir(rootDir), value)
	}
	c.SourceDirectory = value
}

const S3RegionDefault string = "us-east-1"

func (c *Config) loadS3Region(raw configMap) {
	val := cast.ToString(raw.get("s3_region"))
	if val == "" {
		val = S3RegionDefault
	}
	c.S3Region = val
}

func (c *Config) loadS3Bucket(raw configMap) {
	val := cast.ToString(raw.get("s3_bucket"))
	if val == "" {
		log.Panic("Configuration error: s3_bucket is missing")
	}
	c.S3Bucket = val
}

func (c *Config) loadRedirects(raw configMap) {
	if raw.get("redirects") == nil {
		return
	}

	if redirects, ok := raw.get("redirects").(configMap); ok {
		for sourcePath, destPath := range redirects {
			path := cast.ToString(sourcePath)
			dest := cast.ToString(destPath)
			if strings.HasPrefix(path, "/") {
				c.redirects[path[1:]] = dest
			} else {
				c.redirects[path] = dest
			}
		}
	} else {
		log.Panic("Configuration error: redirects should be a list of paths and destinations")
	}
}

func (c *Config) Redirects() map[string]string {
	return c.redirects
}

func (c *Config) loadIgnore(raw configMap) {
	ignore := raw.get("ignore")
	c.ignorePatterns.loadFromStringSlice(cast.ToStringSlice(ignore), true)
	// Add default ignore patterns
	// TODO: add dot prefix pattern
	c.ignorePatterns.add(".DS_Store", true)
	c.ignorePatterns.add(".git", true)
	c.ignorePatterns.add(".svn", true)
}

func (c *Config) ShouldIgnore(path string) bool {
	return cast.ToBool(c.ignorePatterns.get(path))
}

func (c *Config) loadGzip(raw configMap) {
	v := raw.get("gzip")
	if gzDefault, err := cast.ToBoolE(v); err == nil {
		if gzDefault {
			c.gzipPatterns.add(".html", true)
			c.gzipPatterns.add(".txt", true)
			c.gzipPatterns.add(".css", true)
			c.gzipPatterns.add(".js", true)
			c.gzipPatterns.add(".htm", true)
		}
	} else if gzSlice, err := cast.ToStringSliceE(v); err == nil {
		c.gzipPatterns.loadFromStringSlice(gzSlice, true)
	}
}

func (c *Config) ShouldGzip(path string) bool {
	return cast.ToBool(c.gzipPatterns.get(path))
}

func (c *Config) loadMaxAge(raw configMap) {
	v := raw.get("max-age")
	if defaultMaxAge, err := cast.ToIntE(v); err == nil {
		c.maxAgePatterns.defaultValue = defaultMaxAge
	} else if patterns, ok := v.(configMap); ok {
		for pattern, maxAge := range patterns {
			fmt.Println("pattern", pattern, maxAge)
			c.maxAgePatterns.add(cast.ToString(pattern), cast.ToInt(maxAge))
		}
	} else {
		log.Panicf("invalid max-age config: %T", v)
	}
}

func (c *Config) MaxAge(path string) int {
	return cast.ToInt(c.maxAgePatterns.get(path))
}
