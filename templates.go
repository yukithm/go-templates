package templates

import (
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type files struct {
	base     string
	noBase   bool
	partials []string
}

type Templates struct {
	Dir         string
	DefaultBase string
	LeftDelim   string
	RightDelim  string
	Options     []string
	Funcs       template.FuncMap
}

func (t *Templates) Execute(w io.Writer, name string, data interface{}) error {
	tmpl, err := t.Template(name)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

func (t *Templates) Template(name string) (*template.Template, error) {
	f := t.parseName(name)
	if f.base == "" && !f.noBase {
		f.base = t.DefaultBase
	}

	var tmpl *template.Template
	if !f.noBase && f.base != "" {
		bt, err := t.load(nil, f.base)
		if err != nil {
			return nil, err
		}
		tmpl = bt
	}

	for _, partial := range f.partials {
		_, err := t.load(tmpl, partial)
		if err != nil {
			return nil, err
		}
	}

	names := RequiredTemplates(tmpl)
	for _, name := range names {
		c := tmpl.Lookup(name)
		if c == nil {
			t.load(tmpl, name)
		}
	}

	if t.Options != nil && len(t.Options) > 0 {
		tmpl.Option(t.Options...)
	}

	return tmpl, nil
}

func (t *Templates) parseName(name string) *files {
	var f files
	f.partials = make([]string, 0)

	for _, part := range strings.Split(name, ",") {
		part := strings.TrimSpace(part)
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			switch strings.ToLower(kv[0]) {
			case "base":
				f.base = kv[1]
				if kv[1] == "" {
					f.noBase = true
				}
			}
		} else {
			f.partials = append(f.partials, part)
		}
	}

	return &f
}

func (t *Templates) load(tmpl *template.Template, name string) (*template.Template, error) {
	path := filepath.Clean(filepath.Join(t.Dir, name))
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
