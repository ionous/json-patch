package jsonpatch

import (
	"encoding/json"

	"github.com/ionous/errutil"
)

// At makes a serializable target for writing patches in go.
// FIX: restore the example.
func At(parent, field string) Target {
	return Target{Path(parent), field}
}

// Json turns a string into json data ( for writing patches in go. )
// FIX: restore the example.
func Json(s string) (ret interface{}) {
	var data interface{}
	if e := json.Unmarshal([]byte(s), &data); e != nil {
		ret = errutil.New("error reading json data", e)
	} else {
		ret = data
	}
	return
}
