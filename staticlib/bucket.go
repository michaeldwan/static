package staticlib

import (
	"encoding/hex"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Bucket is an S3 bucket
type Bucket struct {
	region string
	name   string
}

type Object struct {
	Key    string
	Digest []byte
	Size   int64
}

func newBucket(cfg Config) Bucket {
	b := Bucket{
		region: cfg.S3Region,
		name:   cfg.S3Bucket,
	}
	return b
}

// WebsiteEndpoint returns the S3 Static Website endpoint host
func (b Bucket) WebsiteEndpoint() string {
	return fmt.Sprintf("%s.s3-website-%s.amazonaws.com", b.name, b.region)
}

func (b Bucket) Scan() <-chan Object {
	out := make(chan Object)
	go func() {
		fmt.Printf("scan bucket")
		defer close(out)
		listRequest := &s3.ListObjectsInput{
			Bucket: aws.String(b.name),
		}

		s3Client.ListObjectsPages(listRequest, func(page *s3.ListObjectsOutput, lastPage bool) bool {
			for _, obj := range page.Contents {
				out <- Object{
					Key:    *obj.Key,
					Digest: etagToDigest(*obj.ETag),
					Size:   *obj.Size,
				}
			}
			return true
		})
	}()
	return out
}

func etagToDigest(etag string) []byte {
	etagWithoutQuotes := string(etag)[1 : len(etag)-1]
	etagBytes, _ := hex.DecodeString(etagWithoutQuotes)
	return etagBytes
}

func (b Bucket) putFile(file *File, simulate bool) error {
	if simulate {
		return nil
	}

	reader := file.Body()
	defer reader.Close()

	input := &s3.PutObjectInput{
		ACL:           aws.String("public-read"),
		Bucket:        aws.String(b.name),
		Body:          reader,
		ContentLength: aws.Long(file.Size()),
		ContentType:   aws.String(file.ContentType()),
		Key:           aws.String(file.Key()),
		CacheControl:  aws.String(file.CacheControl()),
	}
	_, err := s3Client.PutObject(input)
	return err
}

func (b Bucket) deleteKey(key string, simulate bool) error {
	if simulate {
		return nil
	}

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(b.name),
		Key:    aws.String(key),
	}
	_, err := s3Client.DeleteObject(input)
	return err
}
