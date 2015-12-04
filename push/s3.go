package push

import (
	"encoding/hex"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/michaeldwan/static/context"
)

type DestObject struct {
	Key    string
	Digest []byte
	Size   int64
}

func scanBucket(ctx *context.Context) <-chan *DestObject {
	out := make(chan *DestObject)
	go func() {
		defer close(out)
		listRequest := &s3.ListObjectsInput{
			Bucket: aws.String(ctx.S3Bucket),
		}
		ctx.S3Client().ListObjectsPages(listRequest, func(page *s3.ListObjectsOutput, lastPage bool) bool {
			for _, obj := range page.Contents {
				out <- newDestObject(*obj)
			}
			return true
		})
	}()
	return out
}

func newDestObject(obj s3.Object) *DestObject {
	return &DestObject{
		Key:    *obj.Key,
		Digest: etagToDigest(*obj.ETag),
		Size:   *obj.Size,
	}
}

func etagToDigest(etag string) []byte {
	etagWithoutQuotes := string(etag)[1 : len(etag)-1]
	etagBytes, _ := hex.DecodeString(etagWithoutQuotes)
	return etagBytes
}

func putFile(ctx *context.Context, file *File) error {
	if ctx.Flags.DryRun {
		return nil
	}

	reader := file.Body()
	defer reader.Close()

	input := &s3.PutObjectInput{
		ACL:           aws.String("public-read"),
		Bucket:        aws.String(ctx.S3Bucket),
		Body:          reader,
		ContentLength: aws.Long(file.Size()),
		ContentType:   aws.String(file.ContentType()),
		Key:           aws.String(file.Key()),
		CacheControl:  aws.String(file.CacheControl()),
	}
	_, err := ctx.S3Client().PutObject(input)
	return err
}

func deleteKey(ctx *context.Context, key string) error {
	if ctx.Flags.DryRun {
		return nil
	}

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(ctx.S3Bucket),
		Key:    aws.String(key),
	}
	_, err := ctx.S3Client().DeleteObject(input)
	return err
}

func s3WebsiteEndpoint(region, bucket string) string {
	return fmt.Sprintf("%s.s3-website-%s.amazonaws.com", bucket, region)
}
