package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"github.com/stripe/aws-go/aws"
	"github.com/stripe/aws-go/gen/s3"
)

// File represents a local file that will be uploaded
type File struct {
	Path string
}

var files []File

func visit(path string, f os.FileInfo, err error) error {
	if f.IsDir() {
		return nil
	}
	file := File{path}
	files = append(files, file)
	return nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func computeMd5(filePath string) ([]byte, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return result, err
	}

	return hash.Sum(result), nil
}

func main() {
	flag.Parse()
	root := flag.Arg(0)

	bucket := "michaeldwan.com"

	err := filepath.Walk(root, visit)
	check(err)

	fmt.Printf("will upload: %v\n", files)

	creds := aws.Creds(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")
	fmt.Printf("aws: %v\n", creds)

	client := s3.New(creds, "us-east-1", nil)

	listRequest := s3.ListObjectsRequest{
		Bucket: aws.String(bucket),
	}

	resp, err := client.ListObjects(&listRequest)
	check(err)

	index := make(map[string][]byte)

	for _, obj := range resp.Contents {
		fmt.Printf("%v %v\n", *obj.Key, *obj.ETag)

		// S3 api returns the etag wrapped in quotes, remove them
		etagWithoutQuotes := string(*obj.ETag)[1 : len(*obj.ETag)-1]
		etagBytes, err := hex.DecodeString(etagWithoutQuotes)
		check(err)
		fmt.Printf("  bytes: %v\n", etagBytes)
		index[*obj.Key] = etagBytes
	}

	for _, file := range files {
		io, err := os.Open(file.Path)
		check(err)

		defer io.Close()

		fi, err := os.Stat(file.Path)
		check(err)

		fileSize := fi.Size()

		md5, err := computeMd5(file.Path)
		check(err)

		key, err := filepath.Rel(root, file.Path)
		check(err)

		fmt.Printf("%v\n", key)

		if bytes.Equal(index[key], md5) {
			fmt.Printf("  ---> skip\n")
			continue
		}

		mime := mime.TypeByExtension(filepath.Ext(file.Path))

		fmt.Printf("key: %v (%v)\n", key, mime)

		objectreq := s3.PutObjectRequest{
			ACL:           aws.String("public-read"),
			Bucket:        aws.String(bucket),
			Body:          io,
			ContentLength: aws.Long(fileSize),
			ContentType:   aws.String(mime),
			Key:           aws.String(key),
		}
		_, err = client.PutObject(&objectreq)
		check(err)
		fmt.Printf("%s\n", "https://s3.amazonaws.com/"+bucket+"/"+key)
	}
}
