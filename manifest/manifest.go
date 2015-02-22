package manifest

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/gen/s3"
	"github.com/michaeldwan/webmaster/config"
	"github.com/michaeldwan/webmaster/util"
)

// LocalFile is...
type LocalFile struct {
	Exists   bool
	Path     string
	Digest   []byte
	MimeType string
}

// RemoteFile is...
type RemoteFile struct {
	Exists bool
	Digest []byte
}

// Object represents a local or s3 file in the manifest
type Object struct {
	Key        string
	LocalFile  *LocalFile
	RemoteFile *RemoteFile
}

// Manifest maintains an index of all the local files and s3 objects
type Manifest struct {
	objects map[string]*Object
}

// New returns a new Manifest
func New(c *config.Config) *Manifest {
	manifest := &Manifest{
		objects: make(map[string]*Object),
	}

	remote := manifest.scanBucket(c.Bucket)
	local := manifest.scanDirectory(c.SourceDirectory())

	for {
		select {
		case obj, ok := <-local:
			if ok {
				manifest.setObject(obj)
			} else {
				local = nil
			}
		case obj, ok := <-remote:
			if ok {
				manifest.setObject(obj)
			} else {
				remote = nil
			}
		}

		if local == nil && remote == nil {
			break
		}
	}

	return manifest
}

// Objects returns a map of key/Object pairs
func (m *Manifest) Objects() map[string]*Object {
	return m.objects
}

// Scan passes each key/local/remote tuple to the callback
func (m *Manifest) Scan(callbackFunc func(key string, local *LocalFile, remote *RemoteFile)) {
	for key, obj := range m.Objects() {
		local := obj.LocalFile
		if local == nil {
			local = &LocalFile{}
		}

		remote := obj.RemoteFile
		if remote == nil {
			remote = &RemoteFile{}
		}

		callbackFunc(key, local, remote)
	}
}

func (m *Manifest) scanDirectory(sourcePath string) <-chan *Object {
	fmt.Println("ScanDirectory: ", sourcePath)

	// var wg sync.WaitGroup
	objects := make(chan *Object)

	go func() {
		fmt.Println("walking")
		err := filepath.Walk(sourcePath, func(path string, f os.FileInfo, err error) error {

			if !f.Mode().IsRegular() {
				return nil
			}

			// wg.Add(1)

			// go func() {
			// convert the path to an s3 object key
			key, _ := filepath.Rel(sourcePath, path)

			// compute the md5 of the file
			md5, err := digestFilePath(path)
			if err != nil {
				fmt.Printf("Error calculating md5 on %v: %v\n", path, err)
			}

			mime := mime.TypeByExtension(filepath.Ext(path))

			objects <- &Object{
				Key: key,
				LocalFile: &LocalFile{
					Path:     path,
					Digest:   md5,
					MimeType: mime,
				},
			}
			// wg.Done()
			// }()

			return nil
		})

		util.PanicOnError(err)

		// wg.Wait()
		close(objects)
	}()

	return objects
}

func (m *Manifest) scanBucket(bucket string) <-chan *Object {
	objects := make(chan *Object)

	go func() {
		creds := aws.Creds(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")
		client := s3.New(creds, "us-east-1", nil)
		listRequest := s3.ListObjectsRequest{
			Bucket: aws.String(bucket),
		}

		resp, err := client.ListObjects(&listRequest)
		util.PanicOnError(err)

		for _, obj := range resp.Contents {
			fmt.Printf("%v %v\n", *obj.Key, *obj.ETag)

			// S3 api returns the etag wrapped in quotes, remove them
			etagWithoutQuotes := string(*obj.ETag)[1 : len(*obj.ETag)-1]
			etagBytes, err := hex.DecodeString(etagWithoutQuotes)
			util.PanicOnError(err)

			objects <- &Object{
				Key: *obj.Key,
				RemoteFile: &RemoteFile{
					Digest: etagBytes,
				},
			}
		}

		close(objects)
	}()

	return objects
}

func (m *Manifest) setObject(object *Object) {
	existing, ok := m.objects[object.Key]

	if ok {
		if object.LocalFile != nil {
			existing.LocalFile = object.LocalFile
		} else if object.RemoteFile != nil {
			existing.RemoteFile = object.RemoteFile
		}
	} else {
		m.objects[object.Key] = object
	}
}

func digestFilePath(filePath string) ([]byte, error) {
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
