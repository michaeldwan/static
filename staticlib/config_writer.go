package staticlib

import (
	"os"
	"text/template"
)

func WriteConfigFile(path string) error {
	tmplBody := MustAsset("static.yml")
	tmpl := template.Must(template.New("config.yml.tempkate").Parse(string(tmplBody)))
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}
	return tmpl.Execute(f, map[string]string{
		"sourceDirectory": "./",
		"s3Region":        "us-east-1",
		"s3Bucket":        "YOUR-BUCKET",
	})
}
