package staticlib

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	awsCredentials     *credentials.Credentials
	s3Client           *s3.S3
	cfClient           *cloudfront.CloudFront
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSSessionToken    string
)

func ConfigureAWS(cfg Config) {
	awsCredentials = credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.StaticProvider{Value: credentials.Value{
				AccessKeyID:     AWSAccessKeyID,
				SecretAccessKey: AWSSecretAccessKey,
				SessionToken:    AWSSessionToken,
			}},
			&credentials.EnvProvider{},
			&credentials.SharedCredentialsProvider{},
		})

	s3Client = s3.New(&aws.Config{
		Region:      cfg.S3Region,
		Credentials: awsCredentials,
		// LogLevel: 1,
	})

	cfClient = cloudfront.New(&aws.Config{
		Credentials: awsCredentials,
		// LogLevel: 1,
	})
}
