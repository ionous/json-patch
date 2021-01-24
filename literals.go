package jsonpatch

import (
	"encoding/json"

	"github.com/ionous/errutil"
)

// At makes a serializable target ( for writing patches in go. )
// FIX: restore the example.
func At(parent, field string) Target {
	var full Path
	// manually join the parent and field together (ex. for op.Test )
	if len(field) == 0 {
		full = Path(parent)
	} else {
		full = Path(parent + "['" + field + "']")
	}
	return Target{Path(parent), field, full}
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
