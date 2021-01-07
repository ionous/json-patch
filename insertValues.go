package jsonpatch

import (
	"github.com/ionous/errutil"
)

// InsertValues puts the passed values into the targeted objects; the values are not cloned.
// The number of source and target values must match; returns the number of successful replications.
// ( This is normally used via patch commands. )
func InsertValues(to Cursor, field string, vals []interface{}) (ret int, err error) {
	if toCnt, e := to.Resolve(); e != nil {
		err = e
	} else if fromCnt := len(vals); fromCnt != toCnt {
		err = errutil.Fmt("mismatched copy, from %d to %d", fromCnt, toCnt)
	} else {
		for i := 0; i < toCnt; i++ {
			if obj, ok := to.Element(i).(map[string]interface{}); !ok {
				err = errutil.Fmt("expected a slice of objects; got %T", obj)
				break
			} else {
				obj[field] = vals[i]
				ret++
			}
		}
	}
	return
}
