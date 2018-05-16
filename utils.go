package templates

import (
	"html/template"
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

func stripExt(path string) string {
	ext := filepath.Ext(path)
	if ext == "" {
		return path
	}
	return strings.TrimSuffix(path, ext)
}
