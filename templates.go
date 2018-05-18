package templates

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Templates struct {
	LayoutsDir  string
	ContentsDir string
	PartialsDir string
	LeftDelim   string
	RightDelim  string
	Options     []string
	Funcs       template.FuncMap
}

func (t *Templates) Execute(w io.Writer, name string, data interface{}) error {
	layout, contents := t.parseName(name)

	var tmpl *template.Template
	if layout != "" {
		lt, err := t.load(t.LayoutsDir, layout, nil)
		if err != nil {
			return err
		}
		tmpl = lt
	}

	for _, content := range contents {
		_, err := t.load(t.ContentsDir, content, tmpl)
		if err != nil {
			return err
		}
	}

	names := RequiredTemplates(tmpl)
	for _, name := range names {
		c := tmpl.Lookup(name)
		if c == nil {
			t.load(t.PartialsDir, name, tmpl)
		}
	}

	return tmpl.Execute(w, data)
}

func (t *Templates) load(dir, name string, tmpl *template.Template) (*template.Template, error) {
	d := Dir{
		Path: dir,
	}

	path, found := d.FindFile(name)
	if !found {
		return tmpl, fmt.Errorf("Templates: %s not found", filepath.Join(dir, name))
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return tmpl, err
	}

	nt := t.newTemplate(tmpl, name)
	if tmpl == nil {
		tmpl = nt
	}
	if _, err := nt.Parse(string(buf)); err != nil {
		return tmpl, err
	}

	return tmpl, nil
}

func (t *Templates) loadAll(dir string, ignores []string, tmpl *template.Template) (*template.Template, error) {
	d := Dir{
		Path:        dir,
		ExcludeDirs: ignores,
	}

	err := d.Walk(func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		buf, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		name, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		nt := t.newTemplate(tmpl, name)
		if tmpl == nil {
			tmpl = nt
		}
		if _, err := nt.Parse(string(buf)); err != nil {
			return err
		}

		return nil
	})

	return tmpl, err
}

func (t *Templates) newTemplate(tmpl *template.Template, name string) *template.Template {
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

func (t *Templates) parseName(name string) (string, []string) {
	var layout string
	contents := make([]string, 0)
	for _, part := range strings.Split(name, ",") {
		part := strings.TrimSpace(part)
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			switch strings.ToLower(kv[0]) {
			case "layout":
				layout = kv[1]
			}
		} else {
			contents = append(contents, part)
		}
	}

	return layout, contents
}
