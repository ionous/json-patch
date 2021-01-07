package jsonpatch

import (
	"encoding/json"

	"github.com/ionous/errutil"
)

// Clone copies an in-memory json document.
// ( so we don't accidentally share objects which might later need separate, unique, transformations. )
func Clone(doc interface{}) (ret interface{}, err error) {
	if b, e := json.Marshal(doc); e != nil {
		err = errutil.New("error packing during clone", e)
	} else if e := json.Unmarshal(b, &ret); e != nil {
		err = errutil.New("error unpacking during clone", e)
	}
	return
}
