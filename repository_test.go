package blob

import (
	"testing"
)

func TestRepository(t *testing.T) {
	r, err := Open("fixture")
	if err != nil {
		t.Fatalf("failed to open: %v", err)
	}
	head, err := r.RevParse("HEAD")
	if err != nil {
		t.Fatalf("failed to get HEAD: %v", err)
	}
	if head != "71d67f51ef87c87b2248bb82fda2ed5def2a9f6a" {
		t.Errorf("head != 71d67f51ef87c87b2248bb82fda2ed5def2a9f6a (= %v)", head)
	}
}

func TestObject(t *testing.T) {
	r, err := Open("fixture")
	if err != nil {
		t.Fatalf("failed to open: %v", err)
	}
	id, err := r.GetSHA("HEAD")
	if err != nil {
		t.Fatalf("failed to get HEAD: %v", err)
	}

	object, err := r.Object(id)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	size, err := object.Size()
	t.Log("Object size:", size, err)
}
