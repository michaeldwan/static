package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/gen/s3"
	"github.com/codegangsta/cli"
	"github.com/michaeldwan/webmaster/config"
	"github.com/michaeldwan/webmaster/manifest"
	"github.com/michaeldwan/webmaster/util"
)

func main() {
	app := cli.NewApp()
	app.Name = "webmaster"

	app.Commands = []cli.Command{
		{
			Name:      "manifest",
			ShortName: "m",
			Usage:     "Print the manifest",
			Action: func(c *cli.Context) {
				config := config.New()
				fmt.Printf("config %+v\n", config)

				m := manifest.New(config)

				m.Scan(func(key string, local *manifest.LocalFile, remote *manifest.RemoteFile) {
					fmt.Println(key)
					if local != nil {
						fmt.Println("  Local")
						fmt.Println("    path: ", local.Path)
						fmt.Println("    digest: ", local.Digest)
						fmt.Println("    Mime: ", local.MimeType)
					}
					if remote != nil {
						fmt.Println("  Remote")
						fmt.Println("    digest: ", remote.Digest)
					}
				})
				//
				// m.Scan(func(key string, local *manifest.LocalFile, remote *manifest.RemoteFile) {
				// 	fmt.Println(local.Path)
				// 	m, _ := filepath.Match("assets/*css", key)
				// 	if m {
				// 		fmt.Println("  -> yes")
				// 	} else {
				// 		fmt.Println("  -> no")
				// 	}
				// })

			},
		},

		{
			Name:      "push",
			ShortName: "p",
			Usage:     "Push to AWS",
			Action: func(c *cli.Context) {
				config := config.New()
				m := manifest.New(config)
				creds := aws.Creds(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")
				client := s3.New(creds, "us-east-1", nil)

				m.Scan(func(key string, local *manifest.LocalFile, remote *manifest.RemoteFile) {
					fmt.Println(key)
					if bytes.Equal(local.Digest, remote.Digest) {
						fmt.Println("  ---> skip")
						return
					}

					io, err := os.Open(local.Path)
					util.PanicOnError(err)
					fi, err := os.Stat(local.Path)
					util.PanicOnError(err)
					fileSize := fi.Size()

					objectreq := s3.PutObjectRequest{
						ACL:           aws.String("public-read"),
						Bucket:        aws.String(config.Bucket),
						Body:          io,
						ContentLength: aws.Long(fileSize),
						ContentType:   aws.String(local.MimeType),
						Key:           aws.String(key),
					}
					_, err = client.PutObject(&objectreq)
					fmt.Println("  ---> uploaded")

					io.Close()
				})
			},
		},
	}
	app.Run(os.Args)

	// for _, file := range files {
	// 	io, err := os.Open(file.Path)
	// 	check(err)
	//
	// 	defer io.Close()
	//
	// 	fi, err := os.Stat(file.Path)
	// 	check(err)
	//
	// 	fileSize := fi.Size()
	//
	// 	md5, err := computeMd5(file.Path)
	// 	check(err)
	//
	// 	key, err := filepath.Rel(root, file.Path)
	// 	check(err)
	//
	// 	fmt.Printf("%v\n", key)
	//
	// 	if bytes.Equal(index[key], md5) {
	// 		fmt.Printf("  ---> skip\n")
	// 		continue
	// 	}
	//
	// 	mime := mime.TypeByExtension(filepath.Ext(file.Path))
	//
	// 	fmt.Printf("key: %v (%v)\n", key, mime)
	//
	// 	objectreq := s3.PutObjectRequest{
	// 		ACL:           aws.String("public-read"),
	// 		Bucket:        aws.String(bucket),
	// 		Body:          io,
	// 		ContentLength: aws.Long(fileSize),
	// 		ContentType:   aws.String(mime),
	// 		Key:           aws.String(key),
	// 	}
	// 	_, err = client.PutObject(&objectreq)
	// 	check(err)
	// 	fmt.Printf("%s\n", "https://s3.amazonaws.com/"+bucket+"/"+key)
	// }
}
