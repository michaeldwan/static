package push

import (
  "os"
  "sync"
  "path/filepath"
  "mime"
  "fmt"
  "webmaster/context"
	"compress/gzip"
	"crypto/md5"
  "io"
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

func scanDirectory(ctx *context.Context) <-chan *File {
  out := make(chan *File)
  go func() {
    walkFunc := func(path string, fi os.FileInfo, err error) error {
  		if err != nil { return nil }
      // if fi.IsDir() { return nil }
  		if !fi.Mode().IsRegular() { return nil }
  		if fi.Size() == 0 { return nil }
  		if ctx.ShouldIgnore(path) { return nil }
  		key, _ := filepath.Rel(ctx.SourceDirectory, path)
  		out<-newFile(key, path)
  		return nil
  	}
  	filepath.Walk(ctx.SourceDirectory, walkFunc)
    close(out)
  }()
  return out
}

func generateRedirects(ctx *context.Context) <-chan *File {
  out := make(chan *File)
  go func() {
    for sourcePath, destUrl := range ctx.Redirects() {
      out<-newRedirect(ctx.TempFile(), sourcePath, destUrl)
    }
    close(out)
  }()
  return out
}

func setContentType(in <-chan *File) <-chan *File {
  out := make(chan *File)
  go func() {
    for inFile := range in {
      inFile.contentType = mime.TypeByExtension(filepath.Ext(inFile.path))
      out<-inFile
    }
    close(out)
  }()
  return out
}

func setSize(in <-chan *File) <-chan *File {
  out := make(chan *File)
  go func() {
    for inFile := range in {
      func ()  {
        body := inFile.Body()
      	defer body.Close()
      	if fi, err := body.Stat(); err != nil {
      		panic(err)
      	} else {
      		inFile.size = fi.Size()
      	}
        out<-inFile
      }()
    }
    close(out)
  }()
  return out
}

func setCacheControl(ctx *context.Context, in <-chan *File) <-chan *File {
  out := make(chan *File)
  go func() {
    for inFile := range in {
    	maxAge := ctx.MaxAge(inFile.Key())
    	inFile.cacheControl = fmt.Sprintf("public; max-age=%d", maxAge)
      out<-inFile
    }
    close(out)
  }()
  return out
}

func gzipProcessor(ctx *context.Context, in <-chan *File) <-chan *File {
  out := make(chan *File)
  go func() {
    for inFile := range in {
      out<-gzipFile(ctx, inFile)
    }
    close(out)
  }()
  return out
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
  } else {
    return out
  }
}

func digestProcessor(in <-chan *File) <-chan *File {
  out := make(chan *File)
  go func() {
    for inFile := range in {
      func ()  {
        hash := md5.New()
      	inBody := inFile.Body()
      	defer inBody.Close()
      	var result []byte
      	if _, err := io.Copy(hash, inBody); err != nil {
      		panic(err)
      	}
      	inFile.digest = hash.Sum(result)
      	out<-inFile
      }()
    }
    close(out)
  }()
  return out
}
