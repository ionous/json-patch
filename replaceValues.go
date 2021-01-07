package jsonpatch

import (
	"encoding/json"

	"github.com/ionous/errutil"
)

// ReplaceValues (or add) the 'field' of any objects targeted by the 'parent' path with the passed 'value'.
// ( This is normally used via patch commands. )
func ReplaceValues(parent Cursor, field string, msg json.RawMessage) (ret int, err error) {
	if cnt, e := parent.Resolve(); e != nil {
		err = e
	} else {
		for i := 0; i < cnt; i++ {
			if obj, ok := parent.Element(i).(map[string]interface{}); !ok {
				err = errutil.Fmt("expected a slice of objects; got %T", obj)
				break
			} else {
				var newVal interface{}
				if e := json.Unmarshal(msg, &newVal); e != nil {
					err = errutil.New("couldnt read replacement value because", e)
					break
				} else {
					obj[field] = newVal
					ret++
				}
			}
		}
	}
	return
}
