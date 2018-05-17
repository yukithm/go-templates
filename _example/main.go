package main

import (
	"os"

	"github.com/yukithm/go-templates"
)

func main() {
	dynamic()
}

func template() {
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

func dynamic() {
	tmpls := &templates.DynamicTemplates{
		LayoutsDir:  "./layouts",
		ViewsDir:    "./views",
		PartialsDir: "./partials",
		TemplateExt: ".tmpl",
		StripExt:    true,
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
