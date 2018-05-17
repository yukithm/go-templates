package templates

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"text/template/parse"
)

// RequiredTemplates returns template names which called in t.
func RequiredTemplates(t *template.Template) []string {
	names := make([]string, 0)
	for _, node := range t.Tree.Root.Nodes {
		if node.Type() == parse.NodeTemplate {
			tnode := node.(*parse.TemplateNode)
			names = append(names, tnode.Name)
		}
	}
	return names
}

// AssociateTemplate associates templates by template names in base.
func AssociateTemplate(base *template.Template, content *template.Template) (*template.Template, error) {
	names := RequiredTemplates(base)
	for _, name := range names {
		c := content.Lookup(name)
		if c != nil {
			_, err := base.AddParseTree(name, c.Tree)
			if err != nil {
				return nil, err
			}
		}
	}

	return base, nil
}

func hasExt(path, ext string) bool {
	if ext == "" {
		return true
	}

	return strings.HasSuffix(strings.ToLower(path), strings.ToLower(ext))
}

func stripExt(path string) string {
	ext := filepath.Ext(path)
	if ext == "" {
		return path
	}
	return strings.TrimSuffix(path, ext)
}

func containsPath(path string, list []string) bool {
	cpath := filepath.Clean(path)
	for _, p := range list {
		if cpath == filepath.Clean(p) {
			return true
		}
	}
	return false
}

func findFile(dir, name, ext string) (path string, ok bool) {
	path = filepath.Join(dir, name)
	if _, err := os.Stat(path); err == nil {
		return path, true
	}

	if ext != "" {
		path = filepath.Join(dir, name+ext)
		if _, err := os.Stat(path); err == nil {
			return path, true
		}
	}

	return "", false
}
