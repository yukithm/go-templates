package templates

import (
	"os"
	"path/filepath"
)

type Dir struct {
	Path        string
	ExcludeDirs []string
}

func (d *Dir) FindFile(name string) (path string, ok bool) {
	path = filepath.Clean(filepath.Join(d.Path, name))
	if _, err := os.Stat(path); err == nil {
		return path, true
	}

	return "", false
}

func (d *Dir) Walk(wf filepath.WalkFunc) error {
	return filepath.Walk(d.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return wf(path, info, err)
		}
		if info.IsDir() {
			if d.ExcludeDirs != nil && containsPath(path, d.ExcludeDirs) {
				return filepath.SkipDir
			}
			return nil
		}

		return wf(path, info, err)
	})
}
