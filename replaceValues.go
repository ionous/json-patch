package jsonpatch

import (
	"github.com/PaesslerAG/jsonpath"
	"github.com/ionous/errutil"
)

// ReplaceValues (or add) the 'field' of any objects targeted by the 'parent' path with the passed 'value'.
// If value is nil, delete the field instead.
// ( This is normally used via patch commands. )
func ReplaceValues(doc interface{}, parent, field string, value interface{}) (ret int, err error) {
	if tgt, e := jsonpath.Get(parent, doc); e != nil {
		err = e
	} else if els, ok := tgt.([]interface{}); !ok {
		err = errutil.Fmt("unknown target %T", tgt)
	} else {
		for _, el := range els {
			if obj, ok := el.(map[string]interface{}); !ok {
				err = errutil.Fmt("expected a slice of objects; got %T", el)
				break
			} else {
				if value == nil {
					delete(obj, field)
				} else if newVal, e := Clone(value); e != nil {
					err = e
					break
				} else {
					obj[field] = newVal
				}
				ret++
			}
		}
	}
	return
}
