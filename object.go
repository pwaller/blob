package blob

import (
	"bufio"
	"compress/zlib"
	"crypto"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// Entry point for obtaining an object
func (r *Repository) Object(id SHA) (Object, error) {
	lo := &LooseObject{repoRoot: r.Path, sha: id}

	looseExists, err := exists(lo.Path())
	if err != nil {
		return nil, err
	}
	if looseExists {
		return lo, nil
	}

	return nil, fmt.Errorf("object %q not found (packs not yet implemented)", id)
}

type ObjectType int

const (
	ObjectTypeBlob ObjectType = iota
	ObjectTypeTree
	ObjectTypeCommit
)

func (o ObjectType) String() string {
	switch o {
	case ObjectTypeBlob:
		return "blob"
	case ObjectTypeTree:
		return "tree"
	case ObjectTypeCommit:
		return "commit"
	}
	return "unknown"
}

type ObjectMode string

type Object interface {
	SHA() SHA                       // id
	Type() (ObjectType, error)      // blob, tree, commit
	Size() (int64, error)           // unpacked size
	Reader() (io.ReadCloser, error) // content reader, close method checks SHA
}

type ObjectHeader struct {
	Type ObjectType
	Size int64
}

type LooseObject struct {
	repoRoot string
	sha      SHA

	// Written on tree
	*ObjectMode

	// Contained within object
	mu sync.Mutex
	*ObjectHeader
}

func (lo *LooseObject) Path() string {
	str := lo.sha.String()
	return filepath.Join(lo.repoRoot, "objects", str[:2], str[2:])
}

func (lo *LooseObject) SHA() SHA {
	return lo.sha
}

// Read the header bytes out of an object, and cache them.
func (lo *LooseObject) Header() (*ObjectHeader, error) {
	lo.mu.Lock()
	if lo.ObjectHeader != nil {
		lo.mu.Unlock()
		return lo.ObjectHeader, nil
	}
	lo.mu.Unlock()

	r, err := lo.Reader()
	if err != nil {
		return nil, err
	}
	err = r.Close()
	if err != nil {
		return nil, err
	}

	return lo.ObjectHeader, nil
}

// Returns the object mode and true if known, false otherwise.
func (lo *LooseObject) Mode() (ObjectMode, bool) {
	if lo.ObjectMode == nil {
		return ObjectMode(""), false
	}
	return *lo.ObjectMode, true
}

func (lo *LooseObject) Type() (ObjectType, error) {
	h, err := lo.Header()
	if err != nil {
		return 0, err
	}
	return h.Type, nil
}

func (lo *LooseObject) Size() (int64, error) {
	h, err := lo.Header()
	if err != nil {
		return 0, err
	}
	return h.Size, nil
}

// Returns a ReadCloser on the object complete with header.
func (lo *LooseObject) open() (io.ReadCloser, error) {
	fd, err := os.Open(lo.Path())
	if err != nil {
		return nil, err
	}

	rd, err := zlib.NewReader(fd)
	if err != nil {
		fd.Close()
		return nil, err
	}

	return struct {
		*bufio.Reader
		io.Closer
	}{
		bufio.NewReader(rd),
		fd,
	}, nil
}

type RuneReader interface {
	io.Reader
	ReadBytes(byte) ([]byte, error)
	ReadRune() (rune, int, error)
	UnreadRune() error
}

// Read the header out of an object
func (lo *LooseObject) readHeader(r RuneReader) error {
	lo.mu.Lock()
	defer lo.mu.Unlock()

	if lo.ObjectHeader != nil {
		_, err := r.ReadBytes('\x00') // already read, consume and continue
		return err
	}

	h := &ObjectHeader{}
	lo.ObjectHeader = h

	_, err := fmt.Fscanf(r, "%s %d\x00", &h.Type, &h.Size)
	if err != nil {
		return err
	}

	return nil
}

// Returns a reader of an object, which is limited to reading the exact amount
// of header expected. The Close() method also verifies the hash.
func (lo *LooseObject) Reader() (io.ReadCloser, error) {
	r, err := lo.open()
	if err != nil {
		return nil, err
	}

	err = lo.readHeader(r.(RuneReader))
	if err != nil {
		r.Close()
		return nil, err
	}

	r, err = wrapHashVerifier(lo, r)
	if err != nil {
		_ = r.Close() // swallow close error
		return nil, err
	}
	return r, nil
}

type Closer func() error

func (c Closer) Close() error { return c() }

// Wraps `r` in an object which tees the bytes being read into a SHA1, and
// verfies the object on Close(), and also limits the number of bytes being
// read to the object size.
func wrapHashVerifier(o Object, r io.ReadCloser) (io.ReadCloser, error) {
	size, err := o.Size()
	if err != nil {
		return nil, err
	}

	hasher := crypto.SHA1.New()

	teed := io.TeeReader(io.LimitReader(r, size), hasher)

	return struct {
		io.Reader
		io.Closer
	}{
		teed,
		Closer(func() error {

			computed := SHAFromSlice(hasher.Sum(nil))
			if !computed.Equal(o.SHA()) {
				_ = r.Close() // swallow error
				return ErrSHAMismatch{o.SHA(), computed}
			}

			err := r.Close()
			return err
		}),
	}, nil
}
