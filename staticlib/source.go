package staticlib

import (
	"os"
	"path/filepath"
	"sync"
)

type Source struct {
	wg         *sync.WaitGroup
	workingDir workingDir
}

func newSource(cfg Config) Source {
	src := Source{}
	src.workingDir = newWorkingDir(cfg.Path)
	return src
}

func (s Source) clean() error {
	return s.workingDir.clean()
}

func (s Source) process(cfg Config) <-chan File {
	in := make(chan File)
	out := setContentType(in)
	out = setCacheControl(cfg, out)
	out = gzipProcessor(s.workingDir, cfg, out)
	out = digestProcessor(out)
	out = setSize(out)

	go func() {
		s.scanDirectory(cfg, in)
		s.generateRedirects(cfg, in)
		close(in)
	}()

	return out
}

func (s Source) scanDirectory(cfg Config, out chan<- File) {
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
		if cfg.ShouldIgnore(path) {
			return nil
		}
		key, _ := filepath.Rel(cfg.SourceDirectory, path)
		out <- newFile(key, path)
		return nil
	}
	filepath.Walk(cfg.SourceDirectory, walkFunc)
}

func (s Source) generateRedirects(cfg Config, out chan<- File) {
	for sourcePath, destURL := range cfg.Redirects() {
		out <- newRedirect(s.workingDir.tempFile(), sourcePath, destURL)
	}
}
