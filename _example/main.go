package main

import (
	"os"

	"github.com/yukithm/go-templates"
)

func main() {
	tmpls := &templates.Templates{
		LayoutsDir:  "./layouts",
		ViewsDir:    "./views",
		PartialsDir: "./partials",
		TemplateExt: ".tmpl",
		StripExt:    true,
	}
	if err := tmpls.Load(); err != nil {
		panic(err)
	}

	data := struct {
		Name string
	}{
		Name: "Alice",
	}

	if err := tmpls.Execute(os.Stdout, os.Args[1], os.Args[2], data); err != nil {
		panic(err)
	}
}
