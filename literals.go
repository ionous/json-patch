package jsonpatch

import "encoding/json"

// At makes a serializable target for writing patches in go.
// FIX: restore the example.
func At(parent, field string) Target {
	return Target{parent, field}
}

// Json turns a string into json data ( for writing patches in go. )
// FIX: restore the example.
func Json(s string) (ret interface{}) {
	var data interface{}
	if e := json.Unmarshal([]byte(s), &data); e != nil {
		ret = e
	} else {
		ret = data
	}
	return
}
