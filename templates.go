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

// Templates holds layouts, contents and partials templates separately.
type Templates struct {
	LayoutsDir  string
	ViewsDir    string
	PartialsDir string
	TemplateExt string
	StripExt    bool
	LeftDelim   string
	RightDelim  string
	Options     []string
	Funcs       template.FuncMap

	layouts  map[string]*template.Template
	views    map[string]*template.Template
	partials *template.Template
}

func (t *Templates) Load() error {
	if t.LayoutsDir != "" {
		if err := t.LoadLayouts(); err != nil {
			return err
		}
	}

	if t.ViewsDir != "" {
		if err := t.LoadViews(); err != nil {
			return err
		}
	}

	if t.PartialsDir != "" {
		if err := t.LoadPartials(); err != nil {
			return err
		}
	}

	return nil
}

func (t *Templates) AddLayout(name, buf string) error {
	if t.layouts == nil {
		t.layouts = make(map[string]*template.Template)
	}

	tmpl, err := t.newTemplate(nil, name).Parse(buf)
	if err != nil {
		return err
	}
	t.layouts[name] = tmpl

	return nil
}

func (t *Templates) AddLayoutFile(file string) error {
	if !t.isTemplateFile(file) {
		return nil
	}
	return addFile(t.LayoutsDir, file, t.StripExt, t.AddLayout)
}

func (t *Templates) LoadLayouts() error {
	ignores := []string{t.ViewsDir, t.PartialsDir}
	return loadTemplates(t.LayoutsDir, ignores, t.AddLayoutFile)
}

func (t *Templates) AddView(name, buf string) error {
	if t.views == nil {
		t.views = make(map[string]*template.Template)
	}

	tmpl, err := t.newTemplate(nil, name).Parse(buf)
	if err != nil {
		return err
	}
	t.views[name] = tmpl

	return nil
}

func (t *Templates) AddViewFile(file string) error {
	if !t.isTemplateFile(file) {
		return nil
	}
	return addFile(t.ViewsDir, file, t.StripExt, t.AddView)
}

func (t *Templates) LoadViews() error {
	ignores := []string{t.LayoutsDir, t.PartialsDir}
	return loadTemplates(t.ViewsDir, ignores, t.AddViewFile)
}

func (t *Templates) AddPartial(name, buf string) error {
	tmpl, err := t.newTemplate(t.partials, name).Parse(buf)
	if err != nil {
		return err
	}
	if t.partials == nil {
		t.partials = tmpl
	}

	return nil
}

func (t *Templates) AddPartialFile(file string) error {
	if !t.isTemplateFile(file) {
		return nil
	}
	return addFile(t.PartialsDir, file, t.StripExt, t.AddPartial)
}

func (t *Templates) LoadPartials() error {
	ignores := []string{t.LayoutsDir, t.ViewsDir}
	return loadTemplates(t.PartialsDir, ignores, t.AddPartialFile)
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

func addFile(dir, file string, strip bool, addFunc func(name, buf string) error) error {
	path := filepath.Join(dir, file)
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	name := file
	if strip {
		name = stripExt(file)
	}

	return addFunc(name, string(buf))
}

func loadTemplates(dir string, ignoreDirs []string, addFileFunc func(rel string) error) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if ignoreDirs != nil && containsPath(path, ignoreDirs) {
				return filepath.SkipDir
			}
			return nil
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		return addFileFunc(rel)
	})
}

func (t *Templates) isTemplateFile(name string) bool {
	if t.TemplateExt == "" {
		return true
	}

	return strings.HasSuffix(strings.ToLower(name), strings.ToLower(t.TemplateExt))
}

func (t *Templates) Execute(w io.Writer, layout, view string, data interface{}) error {
	lt := t.layouts[layout]
	if lt == nil {
		return fmt.Errorf("Templates: layout not found: %s", layout)
	}
	vt := t.views[view]
	if vt == nil {
		return fmt.Errorf("Templates: view not found: %s", view)
	}

	tmpl, err := lt.Clone()
	if err != nil {
		return err
	}
	if _, err := AssociateTemplate(tmpl, vt); err != nil {
		return err
	}
	if t.partials != nil {
		if _, err := AssociateTemplate(tmpl, t.partials); err != nil {
			return err
		}
	}

	if t.Options != nil {
		tmpl.Option(t.Options...)
	}

	return tmpl.Execute(w, data)
}
