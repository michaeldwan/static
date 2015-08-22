package push

import (
  "webmaster/context"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/awsutil"
  "github.com/aws/aws-sdk-go/service/s3"
  "fmt"
  "log"
  "encoding/hex"
)

func scanBucket(ctx *context.Context) <-chan *DestObject {
	out := make(chan *DestObject)
	go func() {
		defer close(out)
		listRequest := &s3.ListObjectsInput{
			Bucket: aws.String(ctx.S3Bucket),
		}
		resp, err := ctx.S3Client().ListObjects(listRequest)
		if err != nil {
			fmt.Println(awsutil.StringValue(resp))
			log.Fatal(err)
		}
		for _, obj := range resp.Contents {
			out <- newDestObject(*obj)
		}
	}()
	return out
}

func newDestObject(obj s3.Object) *DestObject {
  return &DestObject{
    Key: *obj.Key,
    Digest: etagToDigest(*obj.ETag),
  }
}

func etagToDigest(etag string) []byte {
  etagWithoutQuotes := string(etag)[1 : len(etag)-1]
  etagBytes, _ := hex.DecodeString(etagWithoutQuotes)
  return etagBytes
}

func putFile(ctx *context.Context, file *File) error {
  reader := file.Body()
	defer reader.Close()

  // TODO: retry logic here...
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
  // TODO: retry logic here...
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(ctx.S3Bucket),
		Key:    aws.String(key),
	}
	_, err := ctx.S3Client().DeleteObject(input)
	return err
}
