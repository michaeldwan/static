package staticlib

import (
	"html/template"
	"io"
)

const redirBody string = `<!DOCTYPE html><html><head><meta name="generator" content="static"><meta http-equiv="content-type" content="text/html; charset=utf-8" /><meta http-equiv="refresh" content="0;url={{.}}" /></head></html>`

var redirTemplate = template.Must(template.New("redirect").Parse(redirBody))

func renderRedirect(destURL string, writer io.Writer) error {
	return redirTemplate.Execute(writer, destURL)
}
