package blob

import (
	"testing"
)

func TestRepository(t *testing.T) {
	r, err := Open("fixture")
	if err != nil {
		t.Fatalf("failed to open: %v", err)
	}
	head, err := r.GetRef("HEAD")
	if err != nil {
		t.Fatalf("failed to get HEAD: %v", err)
	}
	if head != "71d67f51ef87c87b2248bb82fda2ed5def2a9f6a" {
		t.Errorf("head != 71d67f51ef87c87b2248bb82fda2ed5def2a9f6a (= %v)", head)
	}
}
