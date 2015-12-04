package context

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Flags struct {
	DryRun      bool
	Force       bool
	Concurrency int
	Verbose     bool

	AWSAccessKeyId     string
	AWSSecretAccessKey string
	AWSSessionToken    string
}

type Context struct {
	Config
	workingDir     workingDir
	awsCredentials *credentials.Credentials
	s3Client       *s3.S3
	cfClient       *cloudfront.CloudFront
	Flags          Flags
}

func New(configPath string) *Context {
	ctx := &Context{}
	ctx.Config = Config{}
	ctx.Config.loadFromPath(configPath)
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
		c.awsCredentials = credentials.NewChainCredentials(
			[]credentials.Provider{
				&credentials.StaticProvider{Value: credentials.Value{
					AccessKeyID:     c.Flags.AWSAccessKeyId,
					SecretAccessKey: c.Flags.AWSSecretAccessKey,
					SessionToken:    c.Flags.AWSSessionToken,
				}},
				&credentials.EnvProvider{},
				&credentials.SharedCredentialsProvider{},
			})
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

func (c *Context) CloudFrontClient() *cloudfront.CloudFront {
	if c.cfClient == nil {
		c.cfClient = cloudfront.New(&aws.Config{
			Credentials: c.AwsCredentials(),
			// LogLevel: 1,
		})
	}
	return c.cfClient
}
