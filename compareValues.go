package jsonpatch

import (
	"encoding/json"

	"github.com/ionous/errutil"
)

// CompareValues compares a path against raw json bytes; FIX: it's currently untested.
// ( This is normally used via patch commands. )
func CompareValues(from Cursor, field string, value json.RawMessage) (retMatches, retMismatch int, err error) {
	if cnt, e := from.Resolve(); e != nil {
		err = e
	} else {
		for i := 0; i < cnt; i++ {
			el := from.Element(i)
			if got, e := json.Marshal(el); e != nil {
				err = errutil.New("compare error compacting el", e)
			} else {
				same := len(got) == len(value) // provisionally
				if same {
					for i, g := range got {
						if value[i] != g {
							same = false
							break
						}
					}
				}
				if same {
					retMatches++
				} else {
					retMismatch++
				}
			}
		}
	}
	return
}
