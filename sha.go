package blob

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

type SHA [20]byte

type SHAer interface {
	SHA() SHA
}

func (s SHA) SHA() SHA {
	return s
}

func (lhs SHA) Equal(rhs SHA) bool {
	return bytes.Equal(lhs[:], rhs[:])
}

func SHAFromSlice(id []byte) SHA {
	return SHA{}
}

func SHAFromString(id string) (SHA, error) {
	var sha SHA

	if hex.DecodedLen(len(id)) != len(sha) {
		return sha, fmt.Errorf("bad SHA %q", id)
	}

	_, err := hex.Decode(sha[:], []byte(id))
	if err != nil {
		return SHA{}, err
	}

	return sha, nil
}

func (s SHA) String() string {
	return hex.EncodeToString(s[:])
}
