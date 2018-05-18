package templates

import (
	"io/ioutil"
	"os"
	"time"
)

type File struct {
	Path     string
	lastLoad time.Time
}

func (f *File) Modified() (bool, error) {
	info, err := os.Stat(f.Path)
	if err != nil {
		return false, err
	}

	return info.ModTime().After(f.lastLoad), nil
}

func (f *File) ReadAll() ([]byte, error) {
	buf, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return nil, err
	}
	f.lastLoad = time.Now()
	return buf, nil
}
