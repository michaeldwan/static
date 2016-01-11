package staticlib

import (
	"compress/gzip"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Source struct {
	path       string
	config     Config
	workingDir workingDir
	contents   map[string]*Content
}

func NewSource(cfg Config) Source {
	src := Source{
		path:   cfg.SourceDirectory,
		config: cfg,
	}
	src.contents = make(map[string]*Content)
	src.workingDir = newWorkingDir(cfg.Path)
	return src
}

func (s *Source) Capture() *Operation {
	op := newOperation()
	go func() {
		defer op.done()
		num := 0
		op.err = filepath.Walk(s.path, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !fi.Mode().IsRegular() {
				return nil
			}
			if fi.Size() == 0 {
				return nil
			}
			if s.config.ShouldIgnore(path) {
				return nil
			}
			content, err := s.newContentFromSourceFile(path, fi)
			if err != nil {
				return err
			}
			s.contents[content.Key] = content
			num++
			op.progressCh <- Progress{Num: num}
			return nil
		})
	}()
	return op
}

func (s *Source) GenerateRedirects() *Operation {
	op := newOperation()
	num := 0
	go func() {
		defer op.done()
		for sourcePath, destURL := range s.config.Redirects() {
			tmp := s.workingDir.tempFile()
			err := renderRedirect(destURL, tmp)
			if err != nil {
				op.err = err
				return
			}
			file, err := newFileRef(tmp.Name())
			if err != nil {
				op.err = err
				return
			}
			content, err := s.newContentFromTempFile(sourcePath, file)
			if err != nil {
				op.err = err
				return
			}
			content.RedirectUrl = destURL
			s.contents[content.Key] = content
			num++
			op.progressCh <- Progress{Num: num}
		}
	}()
	return op
}

func (s *Source) newContentFromSourceFile(path string, fi os.FileInfo) (*Content, error) {
	key := s.keyFromPath(path)
	c := &Content{
		Key:        key,
		ModTime:    fi.ModTime(),
		sourceFile: fileRef{path: path, size: fi.Size()},
	}

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if c.ContentType, err = getContentType(f); err != nil {
		return nil, err
	}
	c.CacheControl = cacheControl(s.config, key)
	return c, err
}

func (s *Source) newContentFromTempFile(key string, tmpFile fileRef) (*Content, error) {
	c := &Content{
		Key:      key,
		ModTime:  time.Now(),
		tempFile: tmpFile,
	}

	f, err := os.Open(tmpFile.path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if c.ContentType, err = getContentType(f); err != nil {
		return nil, err
	}
	c.CacheControl = cacheControl(s.config, key)
	return c, err
}

func (s Source) String() string {
	return fmt.Sprintf("Source: %s", s.path)
}

func (s Source) Clean() error {
	return s.workingDir.clean()
}

func (s *Source) keyFromPath(path string) string {
	key, _ := filepath.Rel(s.path, path)
	return key
}

func (s *Source) ContentForKey(key string) *Content {
	return s.contents[key]
}

func (s *Source) CompressContents() *Operation {
	op := newOperation()
	go func() {
		defer op.done()
		var compressibleContent []*Content
		for _, c := range s.contents {
			if s.config.ShouldGzip(c.ContentType) {
				compressibleContent = append(compressibleContent, c)
			}
		}
		if len(compressibleContent) == 0 {
			return
		}
		progress := Progress{Total: len(compressibleContent)}
		for _, content := range compressibleContent {
			s.compressContent(content)
			progress.Num++
			op.progressCh <- progress
		}
	}()
	return op
}

func (s *Source) compressContent(c *Content) error {
	reader, err := os.Open(c.workingFile().path)
	if err != nil {
		return err
	}
	defer reader.Close()
	writer := s.workingDir.tempFile()
	defer writer.Close()
	compressor := gzip.NewWriter(writer)
	_, err = io.Copy(compressor, reader)
	compressor.Close()
	if err != nil {
		return err
	}
	file, err := newFileRef(writer.Name())
	if err != nil {
		return err
	}
	if c.workingFile().size > file.size {
		c.setTempFile(file)
		c.ContentEncoding = "gzip"
		ratio := float64(file.size) / float64(c.sourceFile.size)
		c.Notes = append(c.Notes, fmt.Sprintf("gzip %.1f%%", ratio*100))
	}

	return nil
}

func (s *Source) DigestContents() *Operation {
	op := newOperation()
	total := len(s.contents)
	num := 0
	go func() {
		defer op.done()
		for _, content := range s.contents {
			file, err := os.Open(content.workingFile().path)
			if err != nil {
				op.err = err
				return
			}
			content.Digest, err = digestReader(file)
			if err != nil {
				op.err = err
				return
			}
			num++
			op.progressCh <- Progress{Num: num, Total: total}
		}
	}()
	return op
}

func getContentType(f *os.File) (string, error) {
	if t := mime.TypeByExtension(filepath.Ext(f.Name())); t != "" {
		return t, nil
	}
	buff := make([]byte, 512)
	if _, err := f.Read(buff); err != nil {
		return "", err
	}
	return http.DetectContentType(buff), nil
}

func cacheControl(cfg Config, key string) string {
	maxAge := cfg.MaxAge(key)
	return fmt.Sprintf("public; max-age=%d", maxAge)
}
