package staticlib

import (
	"html/template"
	"io"
	"os"
)

const redirBody string = `<!DOCTYPE html><html><head><meta http-equiv="content-type" content="text/html; charset=utf-8" /><meta http-equiv="refresh" content="0;url={{.}}" /></head></html>`

func renderRedirect(destUrl string, writer io.Writer) {
	tmpl, err := template.New("redirect").Parse(redirBody)
	if err != nil {
		panic(err)
	}
	if err = tmpl.Execute(writer, destUrl); err != nil {
		panic(err)
	}
}

func newRedirect(writer *os.File, key, destUrl string) File {
	defer writer.Close()
	renderRedirect(destUrl, writer)
	file := newFile(key, writer.Name())
	file.contentType = "text/html"
	file.redirectUrl = destUrl
	return file
}
