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

func (r *Repository) RevParse(name string) (string, error) {
	if name == "HEAD" {
		content, err := r.readFile("HEAD")
		if strings.HasPrefix(content, "ref:") {
			return r.RevParse(strings.TrimSpace(content[len("ref:"):]))
		}
		return content, err
	}

	ref, err := r.readFile(name)
	return strings.TrimSpace(ref), err
}

func (r *Repository) GetSHA(ref string) (SHA, error) {

	ref, err := r.RevParse(ref)
	if err != nil {
		return SHA{}, err
	}
	sha, err := SHAFromString(ref)
	if err != nil {
		return SHA{}, err
	}
	return sha, nil
}
