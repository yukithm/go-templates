package templates

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
)

type Template struct {
	Dir         string
	TemplateExt string
	StripExt    bool
	LeftDelim   string
	RightDelim  string
	Options     []string
	Funcs       template.FuncMap

	files map[string]File
	tmpls map[string]*template.Template
}

func (t *Template) Execute(w io.Writer, name string, data interface{}) error {
	return nil
}

func (t *Template) Lookup(name string) *template.Template {
	if tmpl, ok := t.tmpls[name]; ok {
		return tmpl
	}
	return nil
}

func (t *Template) Load(name string) (*template.Template, error) {
	path, found := findFile(t.Dir, name, t.TemplateExt)
	if !found {
		return nil, fmt.Errorf("Templates: %s not found", filepath.Join(t.Dir, name))
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tmpl := t.newTemplate(nil, name)
	if _, err := tmpl.Parse(string(buf)); err != nil {
		return nil, err
	}

	return tmpl, nil
}

func (t *Template) newTemplate(tmpl *template.Template, name string) *template.Template {
	var nt *template.Template
	if tmpl == nil {
		nt = template.New(name)
	} else {
		nt = tmpl.New(name)
	}

	nt.Delims(t.LeftDelim, t.RightDelim)
	if t.Funcs != nil {
		nt.Funcs(t.Funcs)
	}

	return nt
}
