package jsonpatch

import "encoding/json"

// Clone copies an in-memory json document.
// ( so we don't accidentally share objects which might later need separate, unique, transformations. )
func Clone(doc interface{}) (ret interface{}, err error) {
	var out interface{}
	if b, e := json.Marshal(doc); e != nil {
		err = e
	} else if e := json.Unmarshal(b, &out); e != nil {
		err = e
	}
	return
}
