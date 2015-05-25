package blob

import (
	"fmt"
	"os"
	"path/filepath"
)

type ErrRefNotFound string

func (e ErrRefNotFound) Error() string {
	return fmt.Sprintf("no such ref: %s", string(e))
}

type Repository struct {
	Path string
}

func Open(path string) (*Repository, error) {
	dotGit := filepath.Join(path, ".git")
	_, err := os.Stat(dotGit)
	if err == nil {
		path = dotGit
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	return &Repository{path}, nil
}

func (r *Repository) exists(path string) (bool, error) {
	_, err := os.Stat(filepath.Join(r.Path, path))
	if err == nil {
		return true, nil
	}
	if !os.IsNotExist(err) {
		return false, err
	}
	return false, nil
}
