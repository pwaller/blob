package blob

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

func (r *Repository) readFile(name string) (string, error) {
	content, err := ioutil.ReadFile(filepath.Join(r.Path, name))
	return string(content), err
}

func (r *Repository) RevParse(name string) (Ref, error) {
	if name == "HEAD" {
		content, err := r.readFile("HEAD")
		if strings.HasPrefix(content, "ref:") {
			return r.GetRef(strings.TrimSpace(content[len("ref:"):]))
		}
		return NamedRef(content), err
	}

	ref, err := r.readFile(name)
	return NamedRef(strings.TrimSpace(ref)), err
}

func (r *Repository) GetSHA(ref Ref) (SHA, error) {

	ref, err := r.GetRef(name)
	if err != nil {
		return SHA{}, err
	}
	return
}
