package jsonpatch

import (
	"github.com/ionous/errutil"
)

// ExtractValues of the named fields within the objects selected by the passed cursor.
// If del is true it will remove those elements from the document.
// ( This is normally used via patch commands. )
func ExtractValues(from Cursor, field string, del bool) (ret []interface{}, err error) {
	if cnt, e := from.Resolve(); e != nil {
		err = e
	} else {
		for i := 0; i < cnt; i++ {
			if obj, ok := from.Element(i).(map[string]interface{}); !ok {
				err = errutil.Fmt("expected a slice of objects; got %T", obj)
				break
			} else {
				val := obj[field]
				if del {
					delete(obj, field)
				} else if cloned, e := Clone(val); e != nil {
					err = e
					break
				} else {
					val = cloned
				}
				ret = append(ret, val)
			}
		}
	}
	return
}
