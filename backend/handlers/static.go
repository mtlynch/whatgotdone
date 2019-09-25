package handlers

import (
	"net/http"
	"strings"
)

// staticSystem custom file system handler
type staticSystem struct {
	fs http.FileSystem
}

// Open opens file
func (fs staticSystem) Open(path string) (http.File, error) {
	f, err := fs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := fs.fs.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}
