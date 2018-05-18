package templates

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// DynamicTemplates compiles each call of Execute method.
type DynamicTemplates struct {
	LayoutsDir  string
	ViewsDir    string
	PartialsDir string
	TemplateExt string
	StripExt    bool
	LeftDelim   string
	RightDelim  string
	Options     []string
	Funcs       template.FuncMap
}

func (t *DynamicTemplates) Execute(w io.Writer, layout, view string, data interface{}) error {
	tmpl, err := t.loadTemplate(t.LayoutsDir, layout, nil)
	if err != nil {
		return err
	}

	_, err = t.loadTemplate(t.ViewsDir, view, tmpl)
	if err != nil {
		return err
	}

	if t.PartialsDir != "" {
		ignoreDirs := []string{t.LayoutsDir, t.ViewsDir}
		_, err = t.loadTemplates(t.PartialsDir, ignoreDirs, tmpl)
		if err != nil {
			return err
		}
	}

	if t.Options != nil {
		tmpl.Option(t.Options...)
	}

	return tmpl.Execute(w, data)
}

func (t *DynamicTemplates) loadTemplate(dir, name string, tmpl *template.Template) (*template.Template, error) {
	path, found := findFile(dir, name, t.TemplateExt)
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

func (t *DynamicTemplates) loadTemplates(dir string, ignoreDirs []string, tmpl *template.Template) (*template.Template, error) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if ignoreDirs != nil && containsPath(path, ignoreDirs) {
				return filepath.SkipDir
			}
			return nil
		}

		if t.TemplateExt != "" && !hasExt(path, t.TemplateExt) {
			return nil
		}

		buf, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		name := rel
		if t.StripExt {
			name = stripExt(name)
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

func (t *DynamicTemplates) newTemplate(tmpl *template.Template, name string) *template.Template {
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
