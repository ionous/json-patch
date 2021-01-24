package jsonpatch

import "github.com/ionous/errutil"

// UnknownKey replaces jsonpath's string error for key not found.
type UnknownKey struct {
	og   string
	path string
}

func (e UnknownKey) Error() string {
	return errutil.Sprint(e.og, "in", e.path)
}
