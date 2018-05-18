package main

import (
	"os"

	"github.com/yukithm/go-templates"
)

func main() {
	template()
}

func template() {
	tmpls := &templates.Templates{
		Dir:         "./",
		DefaultBase: "layouts/layout1.tmpl",
	}

	data := struct {
		Title string
		Name  string
	}{
		Title: "Profile",
		Name:  "Alice",
	}

	if err := tmpls.Execute(os.Stdout, os.Args[1], data); err != nil {
		panic(err)
	}
}
