package blob

import (
	"fmt"
)

type ErrSHAMismatch struct {
	Want, Got SHA
}

func (e ErrSHAMismatch) Error() string {
	return fmt.Sprintf("shasum mismatch (want %s != got %s)", e.Want, e.Got)
}
