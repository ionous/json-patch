package jsonpatch

import (
	"bytes"
	"encoding/json"

	"github.com/ionous/errutil"
)

// Clone copies an in-memory json document.
// ( so we don't accidentally share objects which might later need separate, unique, transformations. )
func Clone(doc interface{}) (ret interface{}, err error) {
	if b, e := marshal(doc); e != nil {
		err = errutil.New("error packing during clone", e)
	} else if e := json.Unmarshal(b, &ret); e != nil {
		err = errutil.New("error unpacking during clone", e)
	}
	return
}

// short-cut for encoding without escaping <> in strings.
func marshal(doc interface{}) (ret []byte, err error) {
	var out bytes.Buffer
	js := json.NewEncoder(&out)
	js.SetEscapeHTML(EscapeHTML)
	if e := js.Encode(doc); e != nil {
		err = e
	} else {
		ret = out.Bytes()
	}
	return
}

// EscapeHTML: replace &, >, and < with unicode sequences when copying / comparing strings?
// ( default for jsonpatch is not )
var EscapeHTML = false
