package templates

import "embed"

//go:embed *.tmpl
var FS embed.FS

// GoMainTemplate returns the Go main.go template content
func GoMainTemplate() string {
	data, _ := FS.ReadFile("go_main.tmpl")
	return string(data)
}

// GoModTemplate returns the go.mod template content
func GoModTemplate() string {
	data, _ := FS.ReadFile("go_mod.tmpl")
	return string(data)
}

// HTMLIndexTemplate returns the index.html template content
func HTMLIndexTemplate() string {
	data, _ := FS.ReadFile("html_index.tmpl")
	return string(data)
}
