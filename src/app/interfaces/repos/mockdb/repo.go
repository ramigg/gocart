package mockdb

import "errors"

var (
	errNotFound = errors.New("not found")
)

type Repo struct{}

func (r *Repo) IsNotFoundErr(err error) bool {
	return err == errNotFound
}
