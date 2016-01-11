package staticlib

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Bucket is an S3 bucket
type Bucket struct {
	region  string
	name    string
	objects map[string]*Object
}

type Object struct {
	Key    string
	Digest []byte
	Size   int64
}

func NewBucket(cfg Config) *Bucket {
	b := &Bucket{
		region:  cfg.S3Region,
		name:    cfg.S3Bucket,
		objects: make(map[string]*Object),
	}
	return b
}

// WebsiteEndpoint returns the S3 Static Website endpoint host
func (b Bucket) WebsiteEndpoint() string {
	return fmt.Sprintf("%s.s3-website-%s.amazonaws.com", b.name, b.region)
}

func (b *Bucket) String() string {
	return fmt.Sprintf("Destination: %s (%s)", b.name, b.region)
}

func (b *Bucket) Scan() *Operation {
	op := newOperation()
	go func() {
		defer op.done()
		listRequest := &s3.ListObjectsInput{
			Bucket: aws.String(b.name),
		}

		op.err = s3Client.ListObjectsPages(listRequest, func(page *s3.ListObjectsOutput, lastPage bool) bool {
			for _, obj := range page.Contents {
				object := &Object{
					Key:    *obj.Key,
					Digest: etagToDigest(*obj.ETag),
					Size:   *obj.Size,
				}
				b.objects[object.Key] = object
				op.progressCh <- Progress{Num: len(b.objects)}
			}
			return true
		})
	}()
	return op
}

func etagToDigest(etag string) []byte {
	etagWithoutQuotes := string(etag)[1 : len(etag)-1]
	etagBytes, _ := hex.DecodeString(etagWithoutQuotes)
	return etagBytes
}

func (b Bucket) putFile(content *Content, simulate bool) error {
	if simulate {
		return nil
	}

	reader, err := os.Open(content.workingFile().path)
	if err != nil {
		return err
	}
	defer reader.Close()

	input := &s3.PutObjectInput{
		ACL:             aws.String("public-read"),
		Bucket:          aws.String(b.name),
		Body:            reader,
		ContentLength:   aws.Long(content.Size()),
		ContentType:     aws.String(content.ContentType),
		Key:             aws.String(content.Key),
		CacheControl:    aws.String(content.CacheControl),
		ContentEncoding: aws.String(content.ContentEncoding),
	}
	_, err = s3Client.PutObject(input)
	printAWSError(err)
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
