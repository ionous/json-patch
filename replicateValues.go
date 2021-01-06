package jsonpatch

import (
	"github.com/PaesslerAG/jsonpath"
	"github.com/ionous/errutil"
)

// Replicate copies field(s) from one set of objects to the field(s) of another set.
// The number of source and target values must match; returns the number of successful replications.
// ( This is normally used via patch commands. )
func ReplicateValues(doc interface{}, fromParent, fromField, toParent, toField string) (ret int, err error) {
	if src, e := jsonpath.Get(fromParent, doc); e != nil {
		err = e
	} else if fromEls, ok := src.([]interface{}); !ok {
		err = errutil.Fmt("unknown src %T", src)
	} else if dst, e := jsonpath.Get(toParent, doc); e != nil {
		err = e
	} else if toEls, ok := dst.([]interface{}); !ok {
		err = errutil.Fmt("unknown dst %T", dst)
	} else if fromCnt, toCnt := len(fromEls), len(toEls); fromCnt != toCnt {
		err = errutil.Fmt("mismatched copy, from %d to %d", fromCnt, toCnt)
	} else {
		for i := 0; i < fromCnt; i++ {
			toEl, fromEl := toEls[i], fromEls[i]
			if from, ok := fromEl.(map[string]interface{}); !ok {
				err = errutil.Fmt("expected a slice of objects; got %T", fromEl)
				break
			} else if to, ok := toEl.(map[string]interface{}); !ok {
				err = errutil.Fmt("expected a slice of objects; got %T", toEl)
				break
			} else if newVal, e := Clone(from[fromField]); e != nil {
				err = e
				break
			} else {
				to[toField] = newVal
				ret++
			}
		}
	}
	return
}
