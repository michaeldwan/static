package config

import (
	"os/user"
	"path/filepath"
	"strings"
)

// Config holds configuration data for the site and operation
type Config struct {
	Bucket          string
	sourceDirectory string
}

// New returns an initialized Config object
func New() *Config {
	// creds := aws.Creds(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")
	// client := s3.New(creds, "us-east-1", nil)

	config := &Config{
		Bucket: "michaeldwan.com",
	}
	config.setSourceDirectory("~/src/michaeldwan.com/build/")

	return config
}

func (c *Config) setSourceDirectory(path string) {
	// if path begins with ~/ expand to user's home directory
	if path[:2] == "~/" {
		currentUser, _ := user.Current()
		homeDir := currentUser.HomeDir
		path = strings.Replace(path, "~", homeDir, 1)
	}

	c.sourceDirectory = filepath.Clean(path)
}

// SourceDirectory is the location of the files to upload
func (c *Config) SourceDirectory() string {
	return c.sourceDirectory
}
