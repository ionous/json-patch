package jsonpatch

import (
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/ionous/errutil"
)

// Target a location in a json document
// with a json path pointing to (one or more) objects,
// and a field addressing a member of each objects.
type Target struct {
	Parent Path   `json:"parent"`
	Field  string `json:"field"`
}

// Path is just an indicator of a json path ready string
type Path string

// Select prepares a collection of elements from a json doc pointed to by the path.
func (p Path) Select(doc interface{}) Cursor {
	return Cursor{els: doc, path: string(p), res: nil}
}

// Cursor reads (caches) a collection of objects from a json docs.
// FIX: rather than returning []interface{} provide Resolve() and Element(i) which auto-derefs field
// that will allow transparent handling of nil fields for commands that support that.
type Cursor struct {
	res  error       // nil, cached, or error
	path string      // location of the cursor; for debugging.
	els  interface{} // doc or elements depending
}

const cached = errutil.Error("github.com/ionous/json-patch/cached")

func (c *Cursor) Resolve() (ret int, err error) {
	if els, e := c.resolve(); e != nil {
		err = e
	} else {
		ret = len(els)
	}
	return
}

// Element unpacks a value from a targeted object.
// panics if out of range
func (c *Cursor) Element(i int) (ret interface{}) {
	if els, e := c.resolve(); e != nil {
		panic(e)
	} else {
		ret = els[i]
	}
	return
}

// Resolve returns a collection of objects originally pointed to by a path.
func (c *Cursor) resolve() (ret []interface{}, err error) {
	switch e := c.res; e {
	default:
		err = e
	case cached:
		ret = c.els.([]interface{})
	case nil:
		// FIX: this is surely going to break something at some point
		// Paessler's paths dont handle single quotes by default... but maybe there's a way to make it?
		path := strings.Replace(c.path, "'", `"`, -1)
		if tgt, e := jsonpath.Get(path, c.els); e != nil {
			err = errutil.New("error selecting", c.path, e)
			c.res = err
		} else {
			els, ok := tgt.([]interface{})
			if !ok {
				els = []interface{}{tgt} // convert a single result to a uniform array
			}
			ret, c.els, c.res = els, els, cached
		}
	}
	return
}
