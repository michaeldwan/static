package push

import (
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"sync"

	"github.com/michaeldwan/webmaster/context"
)

func buildSource(ctx *context.Context) <-chan *File {
	in := make(chan *File)
	out := setContentType(in)
	out = setCacheControl(ctx, out)
	out = gzipProcessor(ctx, out)
	out = digestProcessor(out)
	out = setSize(out)

	go func() {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			for f := range scanDirectory(ctx) {
				in <- f
			}
		}()

		go func() {
			defer wg.Done()
			for f := range generateRedirects(ctx) {
				in <- f
			}
		}()

		wg.Wait()
		close(in)
	}()

	return out
}

type pipelineAction func(file *File) *File
type pipelineChan <-chan *File

func pipelineStage(in pipelineChan, f pipelineAction) pipelineChan {
	out := make(chan *File)
	go func() {
		for inFile := range in {
			out <- f(inFile)
		}
		close(out)
	}()
	return out
}

func scanDirectory(ctx *context.Context) pipelineChan {
	out := make(chan *File)
	go func() {
		walkFunc := func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !fi.Mode().IsRegular() {
				return nil
			}
			if fi.Size() == 0 {
				return nil
			}
			if ctx.ShouldIgnore(path) {
				return nil
			}
			key, _ := filepath.Rel(ctx.SourceDirectory, path)
			out <- newFile(key, path)
			return nil
		}
		filepath.Walk(ctx.SourceDirectory, walkFunc)
		close(out)
	}()
	return out
}

func generateRedirects(ctx *context.Context) pipelineChan {
	out := make(chan *File)
	go func() {
		for sourcePath, destURL := range ctx.Redirects() {
			out <- newRedirect(ctx.TempFile(), sourcePath, destURL)
		}
		close(out)
	}()
	return out
}

func setContentType(in pipelineChan) pipelineChan {
	return pipelineStage(in, func(inFile *File) *File {
		inFile.contentType = mime.TypeByExtension(filepath.Ext(inFile.path))
		return inFile
	})
}

func setSize(in pipelineChan) pipelineChan {
	return pipelineStage(in, func(inFile *File) *File {
		body := inFile.Body()
		defer body.Close()
		if fi, err := body.Stat(); err != nil {
			panic(err)
		} else {
			inFile.size = fi.Size()
		}
		return inFile
	})
}

func setCacheControl(ctx *context.Context, in pipelineChan) pipelineChan {
	return pipelineStage(in, func(inFile *File) *File {
		maxAge := ctx.MaxAge(inFile.Key())
		inFile.cacheControl = fmt.Sprintf("public; max-age=%d", maxAge)
		return inFile
	})
}

func gzipProcessor(ctx *context.Context, in pipelineChan) pipelineChan {
	return pipelineStage(in, func(inFile *File) *File {
		return gzipFile(ctx, inFile)
	})
}

func gzipFile(ctx *context.Context, in *File) *File {
	if !ctx.ShouldGzip(in.Key()) {
		return in
	}
	reader := in.Body()
	defer reader.Close()
	writer := ctx.TempFile()
	defer writer.Close()
	compressor := gzip.NewWriter(writer)
	_, err := io.Copy(compressor, reader)
	if err != nil {
		panic(err)
	}
	compressor.Close()
	out := in
	out.setPath(writer.Name())
	out.contentEncoding = "gzip"
	// fmt.Printf("compression saved %%%f (%d, %d)\n", float64(out.Size()) / float64(in.Size()), out.Size(), in.Size())
	if in.Size() <= out.Size() {
		return in
	}
	return out
}

func digestProcessor(in pipelineChan) pipelineChan {
	return pipelineStage(in, func(inFile *File) *File {
		hash := md5.New()
		inBody := inFile.Body()
		defer inBody.Close()
		var result []byte
		if _, err := io.Copy(hash, inBody); err != nil {
			panic(err)
		}
		inFile.digest = hash.Sum(result)
		return inFile
	})
}
