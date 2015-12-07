package staticlib

import (
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io"
	"mime"
	"path/filepath"
)

type pipelineAction func(file File) File
type pipelineChan <-chan File

func pipelineStage(in pipelineChan, f pipelineAction) pipelineChan {
	out := make(chan File)
	go func() {
		for inFile := range in {
			out <- f(inFile)
		}
		close(out)
	}()
	return out
}

func setContentType(in pipelineChan) pipelineChan {
	return pipelineStage(in, func(inFile File) File {
		inFile.contentType = mime.TypeByExtension(filepath.Ext(inFile.path))
		return inFile
	})
}

func setSize(in pipelineChan) pipelineChan {
	return pipelineStage(in, func(inFile File) File {
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

func setCacheControl(cfg Config, in pipelineChan) pipelineChan {
	return pipelineStage(in, func(inFile File) File {
		maxAge := cfg.MaxAge(inFile.Key())
		inFile.cacheControl = fmt.Sprintf("public; max-age=%d", maxAge)
		return inFile
	})
}

func gzipProcessor(workingDir workingDir, cfg Config, in pipelineChan) pipelineChan {
	return pipelineStage(in, func(inFile File) File {
		if cfg.ShouldGzip(inFile.Key()) {
			return gzipFile(workingDir, inFile)
		}
		return inFile
	})
}

func gzipFile(workingDir workingDir, in File) File {
	reader := in.Body()
	defer reader.Close()
	writer := workingDir.tempFile()
	defer writer.Close()
	compressor := gzip.NewWriter(writer)
	_, err := io.Copy(compressor, reader)
	if err != nil {
		panic(err)
	}
	compressor.Close()
	out := in
	out.path = writer.Name()
	out.contentEncoding = "gzip"
	// fmt.Printf("compression saved %%%f (%d, %d)\n", float64(out.Size()) / float64(in.Size()), out.Size(), in.Size())
	if in.Size() <= out.Size() {
		return in
	}
	return out
}

func digestProcessor(in pipelineChan) pipelineChan {
	return pipelineStage(in, func(inFile File) File {
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