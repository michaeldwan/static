package context

import (
  "os"
  "log"
  "github.com/aws/aws-sdk-go/aws/credentials"
	// "github.com/aws/aws-sdk-go/aws/awsutil"
  "github.com/aws/aws-sdk-go/service/s3"
  "github.com/aws/aws-sdk-go/aws"
)

type Context struct {
  *Config
  workingDir workingDir
  awsCredentials  *credentials.Credentials
  s3Client        *s3.S3
}

func New(configPath string) *Context {
  ctx := &Context{}
  ctx.Config = newConfig()
	ctx.Config.loadFile(configPath)
  ctx.workingDir = newWorkingDir(ctx.Path)
  return ctx
}

func (c *Context) TempFile() *os.File {
  return c.workingDir.tempFile()
}

func (c *Context) Clean() {
  c.workingDir.clean()
}

func (c *Context) AwsCredentials() *credentials.Credentials {
	if c.awsCredentials == nil {
		// TODO: detect credentials via other providers before aborting
		creds := credentials.NewEnvCredentials()
		if _, err := creds.Get(); err != nil {
			log.Fatal(err)
		}
		c.awsCredentials = creds
	}
	return c.awsCredentials
}

func (c *Context) S3Client() *s3.S3 {
	if c.s3Client == nil {
		c.s3Client = s3.New(&aws.Config{
		 Region:      c.S3Region,
		 Credentials: c.AwsCredentials(),
		 // LogLevel: 1,
	 })
	}
	return c.s3Client
}
