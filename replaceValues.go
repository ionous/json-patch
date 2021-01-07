package jsonpatch

import (
	"encoding/json"

	"github.com/ionous/errutil"
)

// ReplaceValues adds the passed raw json to the 'field' of objects selected by the passed cursor.
// ( This is normally used via patch commands. )
func ReplaceValues(from Cursor, field string, msg json.RawMessage) (ret int, err error) {
	if cnt, e := from.Resolve(); e != nil {
		err = e
	} else {
		for i := 0; i < cnt; i++ {
			if obj, ok := from.Element(i).(map[string]interface{}); !ok {
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
